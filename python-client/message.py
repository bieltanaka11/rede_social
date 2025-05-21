import json
from typing import Dict, Any


class Message:
    def __init__(self, msg_type: str, timestamp: int, origin: str, payload: Dict[str, Any]):
        self.type = msg_type
        self.timestamp = timestamp
        self.origin = origin
        self.payload = payload

    def to_dict(self):
        return {
            "type": self.type,
            "timestamp": self.timestamp,
            "origin": self.origin,
            "payload": self.payload
        }

    def to_json(self):
        return json.dumps(self.to_dict())

    @staticmethod
    def from_json(data: str):
        try:
            obj = json.loads(data)
            return Message(
                msg_type=obj["type"],
                timestamp=obj["timestamp"],
                origin=obj["origin"],
                payload=obj["payload"]
            )
        except (json.JSONDecodeError, KeyError) as e:
            print(f"[ERRO] Falha ao decodificar mensagem JSON: {e}")
            return None

    def __repr__(self):
        return f"<Message type={self.type} from={self.origin} ts={self.timestamp}>"
