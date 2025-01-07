#!/bin/bash

# Update system packages
sudo apt update -y

# Install MySQL
sudo apt install -y mysql-server

# Start MySQL service
sudo systemctl start mysql
sudo systemctl enable mysql

# Secure installation
sudo mysql_secure_installation

echo "MySQL installation complete!"