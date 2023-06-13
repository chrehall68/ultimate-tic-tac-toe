# Ultimate Tic-Tac-Toe
Coded in Go, AI in Python
## How to Play
Just enter a number that represents the tile you want to claim
as your own. Turns alternate X's to O's. THe red on the console
indicates the current board that your move applies to.

Numbers are converted to tile spaces with row-major order, meaning
0 is top left, 3 is middle left, and 6 is bottom left.

## Compiling buffers
Buffers can be compiled with the following command:
```shell
protoc --go_out=${workspaceRoot} --python_out=${workspaceRoot}/py --proto_path=${workspaceRoot}/proto ${workspaceRoot}/proto/board.proto
```