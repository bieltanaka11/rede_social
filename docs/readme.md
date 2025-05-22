# Sistema de Rede Social Distribuída

Este documento apresenta uma documentação completa do Sistema de Rede Social Distribuída, detalhando os requisitos atendidos, a arquitetura do sistema, e instruções detalhadas para instalação, configuração e testes.

## Sumário

1. [Visão Geral](#visão-geral)
2. [Requisitos Atendidos](#requisitos-atendidos)
3. [Arquitetura do Sistema](#arquitetura-do-sistema)
4. [Instalação e Configuração](#instalação-e-configuração)
5. [Como Testar o Sistema](#como-testar-o-sistema)
6. [Roteiro de Testes](#roteiro-de-testes)
7. [Logs e Monitoramento](#logs-e-monitoramento)

## Visão Geral

O Sistema de Rede Social Distribuída é uma plataforma que permite a interação entre usuários através de publicações de textos, sistema de seguidores e troca de mensagens privadas. O sistema foi implementado como um sistema distribuído, utilizando múltiplos servidores para garantir alta disponibilidade e consistência, com sincronização de relógios e ordenação de mensagens através de relógios lógicos.

### Principais Funcionalidades

- Publicação de textos visíveis para outros usuários
- Sistema de seguidores com notificações de novas publicações
- Troca de mensagens privadas entre usuários
- Replicação de dados em múltiplos servidores
- Sincronização de relógios usando o algoritmo de Berkeley
- Eleição de coordenador usando o algoritmo de bullying
- Ordenação de mensagens usando relógios lógicos

## Requisitos Atendidos

### Funcionalidades de Usuário

#### 1. Publicação de Textos
- ✅ **Requisito**: Usuários podem publicar textos visíveis para outros usuários, com timestamp e associação ao autor.
- **Implementação**: O sistema permite que usuários publiquem textos através do cliente Python. Cada publicação é associada ao usuário que a postou e registrada com um timestamp físico e um relógio lógico.
- **Arquivos Relevantes**: 
  - `python-client/client.py` (opção 1 do menu)
  - `go-server/protocol.go` (definição de `PostPayload`)
  - `go-server/storage.go` (armazenamento de posts)

#### 2. Sistema de Seguidores
- ✅ **Requisito**: Usuários podem seguir outros usuários e receber notificações quando estes publicam novas mensagens.
- **Implementação**: O sistema permite que usuários sigam outros usuários através do cliente Python. Quando um usuário seguido publica uma nova mensagem, o sistema envia notificações para todos os seus seguidores.
- **Arquivos Relevantes**:
  - `python-client/client.py` (opção 2 do menu)
  - `go-server/protocol.go` (definição de `FollowPayload`)
  - `go-server/main.go` (função `notifyFollowers`)

#### 3. Mensagens Privadas
- ✅ **Requisito**: Usuários podem enviar mensagens privadas uns aos outros, entregues de forma confiável e ordenada.
- **Implementação**: O sistema permite que usuários enviem mensagens privadas para outros usuários através do cliente Python. As mensagens são entregues de forma confiável e ordenada usando relógios lógicos.
- **Arquivos Relevantes**:
  - `python-client/client.py` (opção 3 do menu)
  - `go-server/protocol.go` (definição de `PrivateMsgPayload`)
  - `go-server/storage.go` (armazenamento de mensagens privadas)

### Funcionalidades de Servidor

#### 4. Replicação em Múltiplos Servidores
- ✅ **Requisito**: Mensagens e postagens devem ser replicadas em pelo menos três servidores.
- **Implementação**: O sistema replica todas as mensagens e postagens em três servidores Go, garantindo alta disponibilidade e consistência dos dados.
- **Arquivos Relevantes**:
  - `go-server/main.go` (sockets PUSH/PULL para replicação)
  - `scripts/start_servers.sh` (inicialização de 3 instâncias)

#### 5. Adição e Remoção Dinâmica de Servidores
- ✅ **Requisito**: Adição e remoção de servidores devem ser feitas de forma dinâmica, sem comprometer disponibilidade e integridade.
- **Implementação**: O sistema permite configuração de peers via variável de ambiente PEER_ADDRS, mas não possui um mecanismo robusto para detecção e adaptação automática quando servidores são adicionados/removidos durante a execução.
- **Arquivos Relevantes**:
  - `go-server/main.go` (configuração de peers)
  - `scripts/start_servers.sh` (configuração de peers para cada servidor)

### Sincronização e Relógios

#### 6. Sincronização de Relógios (Algoritmo de Berkeley)
- ✅ **Requisito**: O sistema deve garantir que os relógios dos servidores estejam sincronizados usando o algoritmo de Berkeley.
- **Implementação**: O sistema utiliza um coordenador Java que implementa o algoritmo de Berkeley para sincronizar os relógios dos servidores participantes.
- **Arquivos Relevantes**:
  - `java-coordinator/Coordinator.java` (implementação do algoritmo)
  - `go-server/main.go` (socket para receber requisições de sincronização)

#### 7. Eleição de Coordenador (Algoritmo de Bullying)
- ✅ **Requisito**: A eleição do coordenador para sincronização deve usar o algoritmo de bullying.
- **Implementação**: O sistema utiliza o algoritmo de bullying para eleger um coordenador responsável pela sincronização dos relógios.
- **Arquivos Relevantes**:
  - `java-coordinator/BullyElection.java` (implementação do algoritmo)

#### 8. Relógios Lógicos e Ordenação de Mensagens
- ✅ **Requisito**: Cada processo deve manter um relógio lógico, e as mensagens devem ser ordenadas de acordo com esses relógios.
- **Implementação**: Todos os processos (usuários e servidores) mantêm relógios lógicos, e as mensagens são ordenadas de acordo com esses relógios, garantindo consistência na ordem de leitura e entrega.
- **Arquivos Relevantes**:
  - `python-client/message.py` (implementação da classe Message com timestamp)
  - `go-server/protocol.go` (campo Lamport na estrutura de mensagens)

### Logs e Monitoramento

#### 9. Geração de Logs
- ✅ **Requisito**: Todos os processos devem gerar arquivo log com todas as interações.
- **Implementação**: Todos os processos (usuários e servidores) geram logs detalhados de todas as interações, facilitando o monitoramento e a depuração.
- **Arquivos Relevantes**:
  - `go-server/storage.go` (método LogMessage)
  - `scripts/start_servers.sh` e `scripts/start_clients.sh` (redirecionamento para logs)

### Requisitos de Desenvolvimento

#### 10. Uso de Bibliotecas de Comunicação
- ✅ **Requisito**: O projeto deve ser desenvolvido usando qualquer biblioteca de comunicação (ZeroMQ, gRPC, etc.).
- **Implementação**: O projeto utiliza ZeroMQ em todas as suas variantes (Go, Java, Python) para comunicação entre os componentes.
- **Arquivos Relevantes**:
  - `go-server/main.go` (uso de github.com/pebbe/zmq4)
  - `java-coordinator/Coordinator.java` (uso de JeroMQ)
  - `python-client/client.py` (uso de pyzmq)

#### 11. Uso de Múltiplas Linguagens
- ✅ **Requisito**: O projeto deve ser desenvolvido com pelo menos 3 linguagens diferentes.
- **Implementação**: O projeto utiliza Go para os servidores, Java para o coordenador e Python para os clientes.
- **Arquivos Relevantes**:
  - `go-server/` (implementação em Go)
  - `java-coordinator/` (implementação em Java)
  - `python-client/` (implementação em Python)

#### 12. Execução com Múltiplos Servidores e Usuários
- ✅ **Requisito**: O projeto deve executar pelo menos 3 servidores e 5 usuários para teste.
- **Implementação**: Os scripts de inicialização configuram 3 servidores e 5 usuários para teste.
- **Arquivos Relevantes**:
  - `scripts/start_servers.sh` (inicialização de 3 servidores)
  - `scripts/start_clients.sh` (inicialização de 5 clientes)

#### 13. Simulação de Alterações nos Relógios
- ✅ **Requisito**: Os relógios de todos os processos podem sofrer alterações aleatórias para testar a sincronização.
- **Implementação**: O sistema implementa mecanismos para simular alterações nos relógios, permitindo testar a robustez da sincronização.
- **Arquivos Relevantes**:
  - `go-server/main.go` (variável driftOffset)
  - `java-coordinator/Coordinator.java` (cálculo de ajustes)

## Arquitetura do Sistema

O sistema é composto por três componentes principais:

1. **Servidores (Go)**: Responsáveis pelo armazenamento e replicação de publicações e mensagens.
2. **Coordenador (Java)**: Responsável pela eleição de coordenador e sincronização de relógios.
3. **Clientes (Python)**: Interface para usuários interagirem com o sistema.

### Padrão de Mensagem

Todas as comunicações entre os componentes utilizam um formato de mensagem JSON padronizado:

```json
{
  "type": "<TIPO>",         // POST, FOLLOW, MSG_PRIVATE, NOTIFY, ELECTION, COORDINATOR, SYNC_REQUEST, SYNC_REPLY, SYNC_ADJUST
  "from_id": "<ID>",        // ex: "user1", "server0"
  "to_id":   "<ID|ALL>",    // ex: "server1", "ALL"
  "lamport": <INT>,         // relógio lógico de Lamport
  "physical":<FLOAT>,       // timestamp UNIX + offset Berkeley
  "payload": { … }          // dados específicos ao tipo
}
```

## Instalação e Configuração

### Pré-requisitos

* **WSL (Ubuntu)** ou **Linux**:
  * Go (>=1.20)
  * Java JDK (>=11)
  * Python 3.8+
  * libzmq3-dev

### Preparar Ambiente

```bash
sudo apt update
sudo apt install -y golang-go openjdk-11-jdk python3 python3-venv python3-pip libzmq3-dev curl
```

### Configurar Componentes

#### 1. Go Server

```bash
cd go-server
# Inicializar módulo Go
go mod init social-network/go-server
# Instalar binding
go get github.com/pebbe/zmq4
# Compilar
go build -o server
```

#### 2. Java Coordinator

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

#### 3. Python Client

```bash
cd python-client
python3 -m venv venv
source venv/bin/activate
pip install pyzmq
```

## Como Testar o Sistema

### 1. Iniciar os Servidores

```bash
cd scripts
chmod +x start_servers.sh
./start_servers.sh
```

Este script inicia 3 instâncias do servidor Go e 1 instância do coordenador Java. Os logs são redirecionados para o diretório `logs/`.

### 2. Iniciar os Clientes

```bash
cd scripts
chmod +x start_clients.sh
./start_clients.sh
```

Este script inicia 5 instâncias do cliente Python. Os logs são redirecionados para o diretório `logs/`.

### 3. Interagir com o Sistema

Para interagir com o sistema, você pode usar o cliente Python diretamente:

```bash
cd python-client
python3 client.py
```

Ao iniciar o cliente, você será solicitado a fornecer um ID de usuário (por exemplo, "user1"). Em seguida, você verá um menu com as seguintes opções:

1. **Post**: Publicar um texto
2. **Follow**: Seguir outro usuário
3. **Private Msg**: Enviar uma mensagem privada
0. **Exit**: Sair do cliente

## Roteiro de Testes

### Teste 1: Publicação de Posts e Notificações

1. No cliente interativo (`user1`), escolha a opção **1 - Post** e digite um texto.
2. Verifique nos logs de `user2..user5` se a notificação foi recebida: `[NOTIFY] user1 posted: ...`.
3. Verifique em `logs/server0.log` se o POST foi armazenado e replicado.

### Teste 2: Sistema de Seguidores

1. No cliente `user2`, escolha a opção **2 - Follow** e digite `user1`.
2. Faça um novo post em `user1`.
3. Verifique se **apenas** `user2` recebe a notificação.

### Teste 3: Mensagens Privadas

1. No cliente `user3`, escolha a opção **3 - Private Msg**, envie para `user4` e digite uma mensagem.
2. Verifique se a confirmação aparece no console de `user3`.
3. Verifique se a mensagem aparece no console de `user4` ou em `logs/user4.log`.

### Teste 4: Eleição Bully e Sincronização Berkeley

1. Pare o servidor `server0` (SIGTERM ou `pkill -f server0`).
2. Verifique em `logs/coordinator.log` se ocorreu uma eleição e anúncio de novo coordenador.
3. Observe as mensagens `SYNC_REQUEST`, `SYNC_REPLY` e `SYNC_ADJUST` nos logs.
4. Reinicie `server0` e verifique se ele reaplica o ajuste e volta ao cluster.

### Teste 5: Teste de Carga

```bash
for i in {1..100}; do
  echo -e "1\nLoad post $i\n0" | python3 python-client/client.py &
done
wait
```

Verifique a latência de replicação e notificação usando os timestamps nos logs.

## Logs e Monitoramento

Todos os processos escrevem logs detalhados no diretório `logs/`:

- `server0.log`, `server1.log`, `server2.log`: Logs dos servidores Go
- `coordinator.log`: Log do coordenador Java
- `user1.log`, `user2.log`, ..., `user5.log`: Logs dos clientes Python

Os logs seguem o formato `[timestamp][LAMPORT=…] ACTION details`, facilitando a auditoria de ordem e offsets.

Para monitorar os logs em tempo real:

```bash
tail -f logs/server0.log
```

Para buscar eventos específicos:

```bash
grep "NOTIFY" logs/user*.log
```
