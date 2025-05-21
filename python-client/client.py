import zmq
import threading
from message import Message


def notification_listener(sub_socket, lamport_clock):
    while True:
        try:
            raw = sub_socket.recv_string()
            topic, data = raw.split(' ', 1)
            msg = Message.from_json(data)
            lamport_clock[0] = max(lamport_clock[0], msg.timestamp) + 1

            if msg.type == "NOTIFY":
                text = msg.payload.get("text")
                post_id = msg.payload.get("post_id")
                print(f"[NOTIFY][{lamport_clock[0]}] {msg.origin} posted: {text} (post_id={post_id})")
            elif msg.type == "MSG_PRIVATE_DELIVER":
                text = msg.payload.get("text")
                msg_id = msg.payload.get("msg_id")
                print(f"[MSG][{lamport_clock[0]}] Private from {msg.origin}: {text} (msg_id={msg_id})")
        except Exception:
            continue


def main():
    user_id = input("User ID: ").strip()
    lamport_clock = [0]  

    context = zmq.Context()
    req = context.socket(zmq.REQ)
    req.connect("tcp://localhost:5555")

    sub = context.socket(zmq.SUB)
    sub.connect("tcp://localhost:5556")
    sub.setsockopt_string(zmq.SUBSCRIBE, user_id)

    listener = threading.Thread(target=notification_listener, args=(sub, lamport_clock), daemon=True)
    listener.start()

    menu = (
        "\nActions:\n"
        "1 - Post\n"
        "2 - Follow\n"
        "3 - Private Msg\n"
        "0 - Exit\n"
        "Select: "
    )

    while True:
        choice = input(menu).strip()
        if choice == "0":
            break

        lamport_clock[0] += 1

        if choice == "1":
            text = input("Post text: ")
            msg = Message(
                msg_type="POST",
                timestamp=lamport_clock[0],
                origin=user_id,
                payload={"text": text}
            )
        elif choice == "2":
            target = input("Follow user: ")
            msg = Message(
                msg_type="FOLLOW",
                timestamp=lamport_clock[0],
                origin=user_id,
                payload={"followed_id": target}
            )
        elif choice == "3":
            target = input("To user: ")
            text = input("Message: ")
            msg = Message(
                msg_type="MSG_PRIVATE",
                timestamp=lamport_clock[0],
                origin=user_id,
                payload={"msg_id": f"{user_id}-{lamport_clock[0]}", "text": text}
            )
        else:
            print("Invalid option")
            continue

        req.send_string(msg.to_json())
        try:
            reply_raw = req.recv_string()
            reply = Message.from_json(reply_raw)
            lamport_clock[0] = max(lamport_clock[0], reply.timestamp) + 1
            status = reply.payload.get("status")
            if status == 200:
                print(f"[OK][{lamport_clock[0]}] {reply.payload.get('message')}")
            else:
                print(f"[ERR][{lamport_clock[0]}] {reply.payload.get('error')}")
        except Exception as e:
            print(f"No reply or error: {e}")

    req.close()
    sub.close()
    context.term()

if __name__ == "__main__":
    main()
