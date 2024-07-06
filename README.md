# Rate Limiter

Este projeto implementa um rate limiter em Go que pode ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

## Funcionalidades

- Limitação de requisições por endereço IP
- Limitação de requisições por token de acesso
- Configuração de limites e tempos de bloqueio via variáveis de ambiente
- Resposta adequada quando o limite é excedido (HTTP 429)
- Armazenamento de informações de limitação em Redis
- Middleware para integração fácil com servidores web

## Configuração

### Variáveis de Ambiente

O rate limiter pode ser configurado utilizando variáveis de ambiente. Um arquivo `.env` na pasta raiz do projeto pode ser utilizado para definir essas variáveis:

- RATE_LIMIT_IP=10 # Limite de requisições por segundo por IP
- RATE_LIMIT_TOKEN=100 # Limite de requisições por segundo por token
- BLOCK_TIME=5m # Tempo de bloqueio em caso de exceder o limite
- REDIS_ADDRESS=redis:6379 # Endereço do Redis

### Para iniciar o projeto, execute:

- sh
- docker-compose up --build

### Inicialização
#### A função main configura e inicializa o rate limiter:

### Middleware
#### O middleware aplica a limitação de taxa às requisições recebidas. Ele verifica se a requisição excede o limite configurado por IP ou token e bloqueia a requisição se o limite for atingido.

### Armazenamento
#### O armazenamento utiliza Redis para manter o estado do rate limiter. A implementação do armazenamento em Redis está no arquivo redis_storage.go.

## Testes
#### Para garantir a eficácia e a robustez do rate limiter, escrevemos testes automatizados usando o pacote testing do Go. Os testes verificam se o rate limiter funciona conforme esperado para limitação por IP, limitação por token e bloqueio após exceder os limites.

- Arquivo tests/limiter_test.go

Para executar os testes, use o seguinte comando no diretório raiz do seu projeto:

- sh
- go test ./tests -v

## Conclusão
### Este projeto fornece uma solução de rate limiting robusta e configurável em Go, utilizando Redis para armazenamento de estado. A configuração é simples e pode ser realizada via variáveis de ambiente, e a implementação inclui testes automatizados para garantir o funcionamento correto.
