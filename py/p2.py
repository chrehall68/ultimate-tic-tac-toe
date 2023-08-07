# computational libs
import tensorflow as tf
from env import UltimateTicTacToeEnv
import time

MODEL_NAME = "attenppo1.keras"

if __name__ == "__main__":
    env = UltimateTicTacToeEnv()
    model = tf.keras.models.load_model(f"models/{MODEL_NAME}")
    obs = env.reset()
    while 1:
        print("waiting for player to move")
        logits, _ = model(tf.expand_dims(obs, axis=0))
        print(logits)
        action = tf.random.categorical(logits, num_samples=1)[0, 0]
        print(action.numpy())
        obs, reward, done, valid = env.step(action.numpy())
        if done:
            break
    time.sleep(1)
