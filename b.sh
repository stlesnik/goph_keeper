openssl genrsa -out configs/certs/ca.key 2048
openssl req -new -x509 -days 365 -key configs/certs/ca.key -out configs/certs/ca.crt \
  -subj "/C=RU/ST=State/L=City/O=Organization/CN=GophKeeper CA"

# Создаем конфиг для SAN
cat > server.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
IP.1 = 127.0.0.1
EOF

# Генерируем серверный сертификат
openssl genrsa -out configs/certs/server.key 2048
openssl req -new -key configs/certs/server.key -out configs/certs/server.csr \
  -subj "/C=RU/ST=State/L=City/O=Organization/CN=localhost"

openssl x509 -req -in configs/certs/server.csr -CA configs/certs/ca.crt -CAkey configs/certs/ca.key -CAcreateserial \
  -out configs/certs/server.crt -days 365 -extfile server.ext


# Генерируем клиентский сертификат
openssl genrsa -out configs/certs/client.key 2048
openssl req -new -key configs/certs/client.key -out configs/certs/client.csr \
  -subj "/C=RU/ST=State/L=City/O=Organization/CN=GophKeeper Client"

openssl x509 -req -in configs/certs/client.csr -CA configs/certs/ca.crt -CAkey configs/certs/ca.key -CAcreateserial \
  -out configs/certs/client.crt -days 365