#!/bin/sh
# I stash my SSL certs in ~/Sites/ssl
dir=ssl
keyfile=$dir/${1/"*."/star.}.key
cert=$dir/${1/"*."/star.}.crt

openssl req -new -x509 -days 10000 -sha1 -newkey rsa:1024 \
       -nodes -keyout $keyfile -out $cert -subj /O=$1/OU=/CN=$1

echo Generated cert for *.$1
echo
echo Adding cert to your keychain...
security add-trusted-cert $cert
echo
echo "# nginx config"
echo ssl_certificate     $cert\;
echo ssl_certificate_key $keyfile\;
echo
echo "# Apache config"
echo SSLEngine on
echo SSLCertificateFile $cert
echo SSLCertificateKeyFile $keyfile
echo