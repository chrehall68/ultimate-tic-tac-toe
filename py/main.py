import os
import time
import board_pb2


def receive(path: str, tp: type):
    while not os.path.exists(path) or os.stat(path).st_size == 0:
        time.sleep(0.0001)

    ret = tp()
    with open(path, "rb") as file:
        ret.ParseFromString(file.read())
    with open(path, "wb") as file:
        file.truncate()
    return ret


def get_return() -> board_pb2.ReturnMessage:
    return receive("./returnmessage.b", board_pb2.ReturnMessage)


def get_state() -> board_pb2.StateMessage:
    return receive("./statemessage.b", board_pb2.StateMessage)


def make_coord(row, col) -> board_pb2.Coord:
    return board_pb2.Coord(row=row, col=col)


def send_action() -> None:
    print("sending action")
    move = board_pb2.Move(large=make_coord(0, 1), small=make_coord(0, 0))
    action = board_pb2.ActionMessage(move=move)

    print("action is", action)

    with open("./action.b", "wb") as file:
        file.write(action.SerializeToString())
    print("sent action")


if __name__ == "__main__":
    print("original state is:", get_state())
    send_action()
    print("move was valid?:", get_return().valid)
