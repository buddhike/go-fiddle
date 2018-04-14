# Registering and installing certificate

```sh
./generate-certificate.sh

# Install certificate on macOS
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain proxy-ca.pem
```
