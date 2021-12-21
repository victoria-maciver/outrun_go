#!/bin/bash

# Ouput files:
# ca.key: certificate authority private key file (dont share)
# ca.crt: certificate authority trust certificate (dont share with users)
# server.key: server private key, password protected (dont share)
# server.csr: server certificate signing request (shared with ca owner)
# server.crt: server certificate signed by the CA (sent back by CA owner, keep on server)
# server.pem: conversion of server.key into a format gRPC likes (dont share)

# summary:
# private files: ca.key, server.key, server.pem, server.crt
# share files: ca.crt (needed by client), server.csr (needed by CSR)

# Changes these CN's to match your hosts in your environment if needed
SERVER_CN=localhost

# step 1: generate certificate authory + trust certificate (ca.crt)
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN=${SERVER_CN}"

# step 2: generate the server private key (server.key)
openssl genrsa -passout pass:1111 -des3 -out server.key 4096

# step 3: get a certificate signing request from the CA (server.csr)
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=${SERVER_CN}" -addext "subjectAltName = DNS:localhost" 
# step 4: sign the certificate with the CA we created (its called self signing) - server.crt
openssl x509 -req -passin pass:1111 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

# step 5: convert the server certificate to .pem format (server.pem) - usable by gRPC
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem