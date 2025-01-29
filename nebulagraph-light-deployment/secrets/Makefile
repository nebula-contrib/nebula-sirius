# Variables
CA_KEY         = ca.key
CA_CERT        = ca.crt
SERVER_KEY     = server.key
SERVER_CSR     = server.csr
SERVER_CERT    = server.crt
CLIENT_KEY     = client.key
CLIENT_CSR     = client.csr
CLIENT_CERT    = client.crt
PASSWORD_FILE  = password.txt
CONFIG_FILE    = ssl_config.cnf

# Default target
all: server client

# Generate CA
ca: $(CA_KEY) $(CA_CERT)

$(CA_KEY):
	openssl genrsa -out $(CA_KEY) 2048

$(CA_CERT): $(CA_KEY)
	openssl req -x509 -new -nodes -key $(CA_KEY) -sha256 -days 365 -out $(CA_CERT) -subj "/CN=NebulaGraph CA"

# Generate server certificates
# openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256 -extensions v3_req -extfile ssl_config.cnf
server: ca $(SERVER_KEY) $(SERVER_CSR) $(SERVER_CERT)


$(SERVER_KEY):
	openssl genrsa -out $(SERVER_KEY) 2048

$(SERVER_CSR): $(SERVER_KEY)
	openssl req -new -key $(SERVER_KEY) -out $(SERVER_CSR) -subj "/CN=nebula-server" -config $(CONFIG_FILE)

$(SERVER_CERT): $(SERVER_CSR) $(CA_CERT) $(CA_KEY)
	openssl x509 -req -in $(SERVER_CSR) -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial -out $(SERVER_CERT) -days 365 -sha256 -extensions v3_req

# Generate client certificates
client: ca $(CLIENT_KEY) $(CLIENT_CSR) $(CLIENT_CERT)

$(CLIENT_KEY):
	openssl genrsa -out $(CLIENT_KEY) 2048

$(CLIENT_CSR): $(CLIENT_KEY)
	openssl req -new -key $(CLIENT_KEY) -out $(CLIENT_CSR) -subj "/CN=nebula-client" -config $(CONFIG_FILE)

$(CLIENT_CERT): $(CLIENT_CSR) $(CA_CERT) $(CA_KEY)
	openssl x509 -req -in $(CLIENT_CSR) -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial -out $(CLIENT_CERT) -days 365 -sha256

# Generate password file
password:
	echo "your_secure_password" > $(PASSWORD_FILE)

# Clean up generated files
clean:
	rm -f $(CA_KEY) $(CA_CERT) $(SERVER_KEY) $(SERVER_CSR) $(SERVER_CERT) $(CLIENT_KEY) $(CLIENT_CSR) $(CLIENT_CERT) $(PASSWORD_FILE) ca.srl


# Verify the validity of the certificates
verify:
	openssl verify -CAfile $(CA_CERT) $(SERVER_CERT) \
	openssl verify -CAfile $(CA_CERT) $(CLIENT_CERT)

# sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain  ca.crt
.PHONY: all ca server client password clean