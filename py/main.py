# used to send data
import os
import time

# proto definitions
import board_pb2 as pb
import tensorflow as tf

# math
import numpy as np

# constants
ROWS = 3
COLS = 3
CELLS = 9

# exploration parameters
used_valid_actions = dict()  # maps action: how many times it's been used validly
exploration_decay_rate = 0.01  # decay rate for rewarding exploration
exploration_reward = lambda times: np.exp(-exploration_decay_rate * times)

# reward parameters
win_reward = 10
cell_reward = 2


class UltimateTicTacToeEnv:
    obs_dim = (9, 9, 3)
    n_actions = CELLS * CELLS

    def _receive(self, path: str, tp: type):
        while not os.path.exists(path) or os.stat(path).st_size == 0:
            time.sleep(0.0001)

        ret = tp()
        with open(path, "rb") as file:
            ret.ParseFromString(file.read())
        with open(path, "wb") as file:
            file.truncate()
        return ret

    def _get_return(self) -> pb.ReturnMessage:
        return self._receive("./returnmessage.b", pb.ReturnMessage)

    def _get_state(self) -> pb.StateMessage:
        return self._receive("./statemessage.b", pb.StateMessage)

    def _make_coord(self, idx) -> pb.Coord:
        return pb.Coord(row=idx // COLS, col=idx % COLS)

    def _send_action(self, move) -> None:
        action = pb.ActionMessage(move=move)

        with open("./action.b", "wb") as file:
            file.write(action.SerializeToString())

    def _to_idx(self, coord: pb.Coord) -> int:
        return coord.row * COLS + coord.col

    def _to_multi_idx(self, move: pb.Move) -> int:
        return self._to_idx(move.large) * CELLS + self._to_idx(move.small)

    def _process_state(self, state: pb.StateMessage) -> np.ndarray:
        """
        The structure of the state:
        (9, 9, 3)
        Outer 9 represent board cells
        inner 9 represent the cell spaces
        each space has 3 objects:
            space owner (0, 1, 2) representing if the space is claimed or not
            cell owner (0, 1, 2) representing if the cell the space belongs to is claimed or not
            curcellornot (0, 1); 1 if the space belongs to the current cell, 0 if not
        """
        board_state = np.zeros((9, 9, 3))
        for cell_idx in range(len(state.board.cells)):
            for space_idx in range(len(state.board.cells[cell_idx].spaces)):
                board_state[cell_idx, space_idx, 0] = (
                    state.board.cells[cell_idx].spaces[space_idx].val
                )
                board_state[cell_idx, space_idx, 1] = state.cellowners[cell_idx]
                board_state[cell_idx, space_idx, 2] = (
                    1 if self._to_idx(state.board.curCell) == cell_idx else 0
                )

        return board_state

    def _get_exploration_reward(self, action: int, msg: pb.ReturnMessage) -> float:
        if msg.valid:
            if action not in used_valid_actions:
                used_valid_actions[action] = 1
            else:
                used_valid_actions[action] += 1
            return exploration_reward(used_valid_actions[action])
        return 0

    def _get_win_reward(self, msg: pb.ReturnMessage) -> float:
        """
        Get's the reward for winning if the game was won
        """
        # the turn sent in the return message should still be the caller's turn
        if msg.state.winner == msg.state.turn:
            return win_reward
        return 0

    def _get_cell_reward(self, msg: pb.ReturnMessage) -> float:
        """
        Get's the reward for claiming a cell if a cell was claimed
        """
        if self.prev_cellowners == msg.state.cellowners:
            return 0
        elif msg.state.cellowners.count(msg.state.turn) > self.prev_cellowners.count(
            msg.state.turn
        ):
            return cell_reward
        return 0

    def _get_reward(self, action: pb.Move, msg: pb.ReturnMessage) -> float:
        return (
            self._get_exploration_reward(action, msg)
            + self._get_cell_reward(msg)
            + self._get_win_reward(msg)
        )

    # public section
    def observe(self) -> np.ndarray:
        return self._process_state(self._get_state())

    def step(self, action: int) -> bool:
        """
        Returns:
            - done / not done
        """
        self._send_action(self.to_move(action))
        ret_message = self._get_return()

        print("was the move valid?", ret_message.valid)
        done = ret_message.state.done
        return done

    def to_move(self, idx: int) -> pb.Move:
        outer_idx = idx // CELLS
        inner_idx = idx % CELLS

        return pb.Move(
            large=self._make_coord(outer_idx), small=self._make_coord(inner_idx)
        )


if __name__ == "__main__":
    env = UltimateTicTacToeEnv()
    model = tf.keras.models.load_model("player2.keras")
    while 1:
        print("waiting for player to move")
        obs = env.observe()
        logits, value = model(tf.expand_dims(obs, axis=0))
        action = tf.argmax(logits, axis=1)[0]
        print(action.numpy())
        done = env.step(action.numpy())
        if done:
            break
