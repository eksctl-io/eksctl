CA_FILES := ca.csr ca-key.pem ca.pem
CERT_FILES := test-server.csr test-server-key.pem test-server.pem

all: $(CERT_FILES)

clean:
	rm -f $(CERT_FILES) $(CA_FILES)

$(CA_FILES): ca.json
	cfssl gencert \
		-initca ca.json \
	| cfssljson -bare ca

$(CERT_FILES): test-server.json config.json $(CA_FILES)
	cfssl gencert \
		-ca ca.pem \
		-ca-key ca-key.pem \
		-config config.json \
		-profile test-server \
		test-server.json \
	| cfssljson -bare test-server
