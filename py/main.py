# computational libs
import tensorflow as tf
import numpy as np

# used to send data
import os
import time
import socket
import board_pb2 as pb

# misc
from typing import Tuple
import configparser


config = configparser.ConfigParser()
config.read("train.ini")

# model
MODEL_NAME = "ppo15.keras"

# Env Constants
MAX_TIMESTEPS = config["ENV"].getint("MAX_TIMESTEPS")

# board constants
ROWS = config["ENV"].getint("ROWS")
COLS = config["ENV"].getint("COLS")
CELLS = config["ENV"].getint("CELLS")

# socket constants
S_PORT = config["ENV"].getint("S_PORT")
A_PORT = config["ENV"].getint("A_PORT")
R_PORT = config["ENV"].getint("R_PORT")
MAX_MSG_SIZE = config["ENV"].getint("MAX_MSG_SIZE")

# reward parameters
WIN_REWARD = config["REWARD"].getfloat("WIN_REWARD")
CELL_REWARD = config["REWARD"].getfloat("CELL_REWARD")
VALID_REWARD = config["REWARD"].getfloat("VALID_REWARD")
INVALID_PENALTY = config["REWARD"].getfloat("INVALID_PENALTY")
LOSS_PENALTY = config["REWARD"].getfloat("LOSS_PENALTY")

# misc
SLEEP_TIME = config["ENV"].getfloat("SLEEP_TIME")


# env
class UltimateTicTacToeEnv:
    obs_dim = (9, 9, 4)
    n_actions = CELLS * CELLS

    def __init__(self) -> None:
        self.s_conn, self.a_conn, self.r_conn = None, None, None

    def _receive(self, conn: socket.socket, tp: type):
        ret = tp()
        b = conn.recv(MAX_MSG_SIZE)
        ret.ParseFromString(b)
        return ret

    def _get_return(self) -> pb.ReturnMessage:
        return self._receive(self.r_conn, pb.ReturnMessage)

    def _get_state(self) -> pb.StateMessage:
        return self._receive(self.s_conn, pb.StateMessage)

    def _make_coord(self, idx) -> pb.Coord:
        return pb.Coord(row=idx // COLS, col=idx % COLS)

    def _send_action(self, move) -> None:
        action = pb.ActionMessage(move=move)
        self.a_conn.send(action.SerializeToString())

    def _to_idx(self, coord: pb.Coord) -> int:
        return coord.row * COLS + coord.col

    def _to_multi_idx(self, move: pb.Move) -> int:
        return self._to_idx(move.large) * CELLS + self._to_idx(move.small)

    def _process_state(self, state: pb.StateMessage) -> np.ndarray:
        """
        The structure of the state:
        (9, 9, 4)
        Outer 9 represent board cells
        inner 9 represent the cell spaces
        each space has 3 objects:
            space owner (0, 1, 2) representing if the space is claimed or not
            cell owner (0, 1, 2) representing if the cell the space belongs to is claimed or not
            curcellornot (0, 1); 1 if the space belongs to the current cell, 0 if not
            turn (1, 2) 1 if the current turn is player1, 2 if the current turn is player2
        """
        board_state = np.zeros(self.obs_dim)
        for cell_idx in range(len(state.board.cells)):
            for space_idx in range(len(state.board.cells[cell_idx].spaces)):
                board_state[cell_idx, space_idx, 0] = (
                    state.board.cells[cell_idx].spaces[space_idx].val
                )
                board_state[cell_idx, space_idx, 1] = state.cellowners[cell_idx]
                board_state[cell_idx, space_idx, 2] = (
                    1 if self._to_idx(state.board.curCell) == cell_idx else 0
                )
                board_state[cell_idx, space_idx, 3] = state.turn

        return board_state

    def _get_exploration_reward(self, action: int, msg: pb.ReturnMessage) -> float:
        if msg.valid:
            return VALID_REWARD
        return INVALID_PENALTY

    def _get_win_reward(self, msg: pb.ReturnMessage) -> float:
        """
        Get's the reward for winning if the game was won
        """
        # the turn sent in the return message should still be the caller's turn
        if msg.state.winner == msg.state.turn:
            if self.player_turn:
                self.won = True
                return WIN_REWARD
            else:
                self.lost = True
                return LOSS_PENALTY
        return 0

    def _get_cell_reward(self, msg: pb.ReturnMessage) -> float:
        """
        Get's the reward for claiming a cell if a cell was claimed
        """
        if self.prev_cellowners == msg.state.cellowners:
            return 0
        elif list(msg.state.cellowners).count(
            msg.state.turn
        ) > self.prev_cellowners.count(msg.state.turn):
            self.prev_cellowners = list(msg.state.cellowners)
            return CELL_REWARD
        return 0

    def _get_reward(self, action: pb.Move, msg: pb.ReturnMessage) -> float:
        return (
            self._get_exploration_reward(action, msg)
            + self._get_cell_reward(msg)
            + self._get_win_reward(msg)
        )

    def _step(self, action: int) -> Tuple[np.ndarray, float, bool, bool]:
        """
        Updates self.done
        """
        # send action and get response
        self._send_action(self.to_move(action))
        ret_message = self._get_return()

        # return information
        reward = self._get_reward(action, ret_message)
        self.done = ret_message.state.done
        if not self.done:
            return self.observe(), reward, self.done, ret_message.valid
        else:
            return (
                self._process_state(ret_message.state),
                reward,
                self.done,
                ret_message.valid,
            )

    def _reset_vars(self):
        self.prev_cellowners = [pb.NONE] * 9
        self.cur_state = None  # the current state; used for debugging
        self.won = False  # whether or not the player won
        self.done = False  # if the game is over
        self.lost = False  # whether or not the player lost
        self.player_turn = True  # whether or not it is the player's turn
        self.cur_timestep = 0

        self.s_conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.a_conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.r_conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

        self.s_conn.connect(("", 8000))
        self.a_conn.connect(("", 8001))
        self.r_conn.connect(("", 8002))

    # public section
    def observe(self) -> np.ndarray:
        """
        Updates self.cur_state and self._turn
        """
        state = self._get_state()
        self._turn = state.turn
        self.cur_state = self._process_state(state)
        return self.cur_state

    def step(self, action: int) -> Tuple[np.ndarray, float, bool, bool]:
        """
        Updates current timestep

        Returns:
            - next state
            - reward for the action
            - done / not done
            - valid / invalid
        """
        self.player_turn = True
        obs, reward, done, valid = self._step(action)
        self.cur_timestep += 1

        return obs, reward, done, valid

    def turn(self):
        return self._turn

    def reset(self) -> np.ndarray:
        self._reset_vars()
        return self.observe()

    def cleanup(self):
        os.system("killall -q uttt")
        if self.s_conn is not None:
            self.s_conn.close()
            self.r_conn.close()
            self.a_conn.close()

    def __del__(self):
        self.cleanup()

    def to_move(self, idx: int) -> pb.Move:
        outer_idx = idx // CELLS
        inner_idx = idx % CELLS

        return pb.Move(
            large=self._make_coord(outer_idx), small=self._make_coord(inner_idx)
        )


if __name__ == "__main__":
    env = UltimateTicTacToeEnv()
    model = tf.keras.models.load_model(f"models/{MODEL_NAME}")
    obs = env.reset()
    while 1:
        print("waiting for player to move")
        logits, value = model(tf.expand_dims(obs, axis=0))
        action = tf.argmax(logits, axis=1)[0]
        print(action.numpy())
        obs, reward, done, valid = env.step(action.numpy())
        if done:
            break
    time.sleep(1)
