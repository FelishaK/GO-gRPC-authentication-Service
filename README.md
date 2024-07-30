# About the project
## This a simple gRPC authentication service

### Features implemented: Register, login, retrieve inforamtion about authenticated user, refresh stale tokens

## Core stack: GO, MongoDB

# How to user the project
 1. Clone the repo: git clone https://github.com/FelishaK/GO-gRPC-authentication-Service.git
 2. mkdir config, cd config
 3. touch config.yaml and put this in there:
  ```yaml
  env: "local"
  access_token_ttl: 1h
  refresh_token_ttl: 720h
  database:
    mongo_user: "mongouser"
    mongo_password: "Password"
    mongo_hostname: "mongo" 
    mongo_port: 27017
    mongo_db_name: "auth"
    timeout: 1s
  grpc:
    port: 20005
    grpc_host: "0.0.0.0"
    timeout: 4s
```
  4. touch private_key.pem public_key.pem
  5. echo "-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgG7e8JC/nQjo4yUuanft38GbpltbdbZ2STVtNNk+gS367+seCsED
zvIQKAibf/vAamVtEdr820TznQhw+WkYTKhcTaLw33RzSPDf/z2wPja5O5P1kjXR
Sk0etrVSYWDoByp1pdtCjoOBOo8L7NDEMEHZ1qoabuvLffrLM1r3JhrtAgMBAAEC
gYAEm7ecDJrWV/e4/+jk+zolrfaILZEC+H+qfNOJhBOSea+nMiR4SVQ8s3c2hGAZ
crH5bUMkuwXSI94PD8MOHzhwaP5HzZKcdnMNHvXsqx+zEX1s3kAeRRHQCYiq8s9C
ZGCjkZ91sOK5yW1d327Ar+3QlsFRxlLTnw7yfcUm1HowAQJBALKmv24WDsPCB/Pu
zszXwH28XWuNMf76Tf58uKElOBpvQy2vfq1nuVdreUwLX53XN7qrN+4px9L8Ip2g
3ulFre0CQQCe35PUgoiL3EiXkokaoCWsmg646qbbs+L+U+FMqE+Oh754K850CBbW
Dw6IL58iurq/4hecPLACyqB5ioiurIEBAkAalsrC/bFw3T4FxjMtNadGj3Rv/3HD
e0mEaNep1DpHZOvgrs/xyxBAvJQvBzpR6ag3tif64GkHM9OLFlhW67H5AkBtW2pP
chZ5ZwTUyHn1SN0F5PlTUbnPKxCJjcVcVdKFQmzaHRU8C0Fk0PJozZbVegEICaHE
2oUxNralUrVovrcBAkEAl1adULyB2ongooxzaCZFCf6NpwtdHnT3ZHqJkR4+KmB/
czyaIc6/k9TOjHWRWyltc13Qw5RgjkzIjsLbD5+bRg==
-----END RSA PRIVATE KEY-----"  >> private_key.pem
6. echo "-----BEGIN PUBLIC KEY-----
MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgG7e8JC/nQjo4yUuanft38Gbpltb
dbZ2STVtNNk+gS367+seCsEDzvIQKAibf/vAamVtEdr820TznQhw+WkYTKhcTaLw
33RzSPDf/z2wPja5O5P1kjXRSk0etrVSYWDoByp1pdtCjoOBOo8L7NDEMEHZ1qoa
buvLffrLM1r3JhrtAgMBAAE=
-----END PUBLIC KEY-----" >> public_key.pem

7. cd project
8. make up (may cause an error on windows, if it is, install git bash and repeat or run docker-compose up --build)
