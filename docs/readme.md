# Social Network Distributed — Documentação para Desenvolvedores e Testers

Este documento reúne **tudo** que você precisa saber para **entender**, **compilar**, **executar** e **testar** o projeto de Rede Social Distribuída, implementado com múltiplas linguagens e bibliotecas, incluindo ZeroMQ, Jeromq, Go, Java e Python.

---

## Estrutura Geral do Projeto

```
/social-network
├── go-server/                 # Servidor de Posts e Mensagens (Go + ZeroMQ)
│   ├── main.go                # Ponto de entrada (REP, PUB, PUSH/PULL)
│   ├── storage.go             # Armazenamento em memória e logs
│   └── protocol.go            # Definições de Message e payloads
│   ├── go.mod                 # Módulo Go
│   └── server.exe             # Binário gerado (após go build)
│
├── java-coordinator/          # Coordenação de relógios (Java + Jeromq + Jackson)
│   ├── BullyElection.java     # Eleição Bully via PUB/SUB
│   ├── Coordinator.java       # Algoritmo de Berkeley via REQ/REP
│   └── libs/                  # JARs: jeromq, jackson-core, jackson-databind, jackson-annotations
│       ├── jeromq-0.5.2.jar
│       └── jackson-*.jar
│
├── python-client/             # Cliente CLI (Python + ZeroMQ)
│   ├── client.py              # Menu interativo (REQ + SUB)
│   ├── message.py             # Classe Message JSON + Lamport
│   └── venv/                  # Virtualenv Python (após criação)
│
├── scripts/                   # Scripts de orquestração (bash / PowerShell)
│   ├── start_servers.sh       # Sobe Go Servers e Java Coordinator
│   ├── start_clients.sh       # Sobe 5 instâncias de Python Client
│   ├── start_servers.ps1      # Versão PowerShell (Windows)
│   └── start_clients.ps1      # Versão PowerShell
│
└── logs/                      # Diretório onde os processos escrevem seus logs
    ├── server0.log
    ├── server1.log
    ├── server2.log
    ├── coordinator.log
    ├── user1.log … user5.log
```

---

## Tecnologias e Padrões Usados

| Componente              | Linguagem | Biblioteca / Protocolo  | Padrão                              |
| ----------------------- | --------- | ----------------------- | ----------------------------------- |
| Servidor de Posts       | Go        | github.com/pebbe/zmq4   | REP/REQ, PUB/SUB, PUSH/PULL         |
| Coordenação de Relógios | Java      | Jeromq + Jackson        | PUB/SUB (Bully), REQ/REP (Berkeley) |
| Cliente CLI             | Python    | pyzmq                   | REQ/REP, SUB (notifications)        |
| Logs e Timestamps       | —         | JSON + Lamport + físico | Arquivos texto por processo         |

**Formato de Mensagem (JSON)**

```jsonc
{
  "type": "<TIPO>",         // POST, FOLLOW, MSG_PRIVATE, NOTIFY, ELECTION, COORDINATOR, SYNC_REQUEST, SYNC_REPLY, SYNC_ADJUST
  "from_id": "<ID>",        // ex: "user1", "server0"
  "to_id":   "<ID|ALL>",    // ex: "server1", "ALL"
  "lamport": <INT>,           // relógio lógico de Lamport
  "physical":<FLOAT>,         // timestamp UNIX + offset Berkeley
  "payload": { … }            // dados específicos ao tipo
}
```

---

##  Como Configurar e Executar

### 1. Pré-requisitos

* **WSL (Ubuntu)** ou **Linux**:

  * Go (>=1.20)
  * Java JDK (>=11)
  * Python 3.8+
  * libzmq3-dev
* **Windows** com Git Bash ou PowerShell:

  * Go MSI
  * Java JDK
  * Python 3
  * libzmq (via Chocolatey)

### 2. Preparar Ambiente (WSL)

```bash
sudo apt update
sudo apt install -y golang-go openjdk-11-jdk python3 python3-venv python3-pip libzmq3-dev curl
```

### 3. Go Server

```bash
cd go-server
# Inicializar módulo Go
go mod init social-network/go-server
# Instalar binding
go get github.com/pebbe/zmq4
# Compilar
go build -o server
```

### 4. Java Coordinator

```bash
cd java-coordinator
# Baixar JARs (se ainda não tiver feito)
cd libs
curl -LO https://repo1.maven.org/maven2/org/zeromq/jeromq/0.5.2/jeromq-0.5.2.jar
curl -LO https://repo1.maven.org/maven2/com/fasterxml/jackson/core/jackson-core/2.13.3/jackson-core-2.13.3.jar
curl -LO https://repo1.maven.org/maven2/com/fasterxml/jackson/core/jackson-databind/2.13.3/jackson-databind-2.13.3.jar
curl -LO https://repo1.maven.org/maven2/com/fasterxml/jackson/core/jackson-annotations/2.13.3/jackson-annotations-2.13.3.jar
cd ..
# Compilar
CP=".:libs/jeromq-0.5.2.jar:libs/jackson-core-2.13.3.jar:libs/jackson-databind-2.13.3.jar:libs/jackson-annotations-2.13.3.jar"
javac -cp "$CP" BullyElection.java Coordinator.java
```

### 5. Python Client

```bash
cd python-client
python3 -m venv venv
source venv/bin/activate
pip install pyzmq
```

### 6. Scripts de Orquestração (WSL)

```bash
cd /mnt/c/Users/User/Documents/FEI/Sistemas_distribuidos/sistema_rede_social
chmod +x scripts/*.sh
./scripts/start_servers.sh  # sobe 3 Go + 1 Java coord
./scripts/start_clients.sh  # sobe 5 clients Python
```

---

## 🧪 Roteiro de Testes (Detalhado)

### A. Publicação de Posts & Notificações

1. No client interativo (`user1`), escolha opção **1 - Post** e digite um texto.
2. No log de `user2..user5`, espere pela notificação: `[NOTIFY] user1 posted: ...`.
3. Verifique em `logs/server0.log` que o POST foi armazenado e replicado.

### B. Seguidores & Filtragem

1. No `user2`, escolha **2 - Follow**, digite `user1`.
2. Faça novo post em `user1`.
3. Garanta que **apenas** `user2` recebe NOTIFY.

### C. Mensagens Privadas

1. Em `user3`, opção **3 - Private Msg**, envie para `user4`.
2. Confira ACK no console/`logs/user3.log`.
3. Confira mensagem em `user4` ou `logs/user4.log`.

### D. Eleição Bully & Sincronização Berkeley

1. Pare o `server0` (SIGTERM ou `pkill -f server0`).
2. Em `logs/coordinator.log`, veja eleição e anúncio de novo coordenador.
3. Observe `SYNC_REQUEST`, múltiplos `SYNC_REPLY` e `SYNC_ADJUST` enviados.
4. Reinicie `server0`, confira que ele reaplica o ajuste e volta ao cluster.

### E. Teste de Carga

```bash
for i in {1..100}; do
  echo -e "1\nLoad post $i\n0" | python3 python-client/client.py &
done
wait
```

* Meça latência de replicação e notificação com timestamps de log.

---

## 📄 Logs e Depuração

* Todos os processos escrevem padrões `[timestamp][LAMPORT=…] ACTION details` em `logs/`.
* Use `grep`, `tail -f` e ferramentas de parsing para auditar ordens e offsets.

---

* Gabriel Carvalho
* RA: 22.121.112-1
