#!/bin/bash
# Willpower Forge - Service Installation Script
# This script automates the installation of Willpower Forge as a systemd service

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="willpower-forge"
INSTALL_DIR="/opt/${APP_NAME}"
DATA_DIR="/var/lib/${APP_NAME}"
SERVICE_FILE="${APP_NAME}.service"
EXECUTABLE="willpower-forge-linux"

echo -e "${GREEN}=====================================${NC}"
echo -e "${GREEN}Willpower Forge Service Installer${NC}"
echo -e "${GREEN}=====================================${NC}"
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}Error: This script must be run as root (use sudo)${NC}"
   exit 1
fi

# Check if executable exists
if [ ! -f "./${EXECUTABLE}" ]; then
    echo -e "${RED}Error: ${EXECUTABLE} not found in current directory${NC}"
    echo "Please make sure you're running this script from the project root"
    exit 1
fi

# Check if service file exists
if [ ! -f "./${SERVICE_FILE}" ]; then
    echo -e "${RED}Error: ${SERVICE_FILE} not found in current directory${NC}"
    exit 1
fi

# Get the current user (the one who ran sudo)
ACTUAL_USER=${SUDO_USER:-$USER}
ACTUAL_GROUP=$(id -gn $ACTUAL_USER)

echo -e "${YELLOW}Configuration:${NC}"
echo "  User: $ACTUAL_USER"
echo "  Group: $ACTUAL_GROUP"
echo "  Install directory: $INSTALL_DIR"
echo "  Data directory: $DATA_DIR"
echo ""
read -p "Continue with installation? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Installation cancelled."
    exit 1
fi

# Step 1: Create directories
echo -e "${GREEN}Step 1/6: Creating directories...${NC}"
mkdir -p "$INSTALL_DIR"
mkdir -p "$DATA_DIR"

# Step 2: Copy executable
echo -e "${GREEN}Step 2/6: Copying executable...${NC}"
cp "./${EXECUTABLE}" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/${EXECUTABLE}"

# Step 3: Set permissions
echo -e "${GREEN}Step 3/6: Setting permissions...${NC}"
chown -R $ACTUAL_USER:$ACTUAL_GROUP "$DATA_DIR"

# Step 4: Create service file with correct user
echo -e "${GREEN}Step 4/6: Creating service file...${NC}"
cat > "/etc/systemd/system/${SERVICE_FILE}" <<EOF
[Unit]
Description=Willpower Forge - Goal Tracking Application
Documentation=https://github.com/ShengWang1017/willpower
After=network.target

[Service]
Type=simple
User=$ACTUAL_USER
Group=$ACTUAL_GROUP
WorkingDirectory=$DATA_DIR
ExecStart=$INSTALL_DIR/$EXECUTABLE
Restart=no

# Security settings
NoNewPrivileges=true
PrivateTmp=true

# Environment variables
# Set PORT to change the listening port (default: 5173)
Environment="PORT=5173"

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=$APP_NAME

[Install]
WantedBy=multi-user.target
EOF

# Step 5: Reload systemd
echo -e "${GREEN}Step 5/6: Reloading systemd...${NC}"
systemctl daemon-reload

# Step 6: Enable and start service
echo -e "${GREEN}Step 6/6: Enabling and starting service...${NC}"
systemctl enable $APP_NAME
systemctl start $APP_NAME

# Wait a moment for service to start
sleep 2

# Check service status
if systemctl is-active --quiet $APP_NAME; then
    echo ""
    echo -e "${GREEN}=====================================${NC}"
    echo -e "${GREEN}Installation completed successfully!${NC}"
    echo -e "${GREEN}=====================================${NC}"
    echo ""
    echo -e "${GREEN}Service Status:${NC}"
    systemctl status $APP_NAME --no-pager -l
    echo ""
    echo -e "${GREEN}Useful Commands:${NC}"
    echo "  Start service:   sudo systemctl start $APP_NAME"
    echo "  Stop service:    sudo systemctl stop $APP_NAME"
    echo "  Restart service: sudo systemctl restart $APP_NAME"
    echo "  View status:     sudo systemctl status $APP_NAME"
    echo "  View logs:       sudo journalctl -u $APP_NAME -f"
    echo ""
    echo -e "${GREEN}Access the application:${NC}"
    echo "  URL: http://localhost:8080"
    echo ""
else
    echo ""
    echo -e "${RED}=====================================${NC}"
    echo -e "${RED}Installation completed with errors!${NC}"
    echo -e "${RED}=====================================${NC}"
    echo ""
    echo -e "${YELLOW}Service failed to start. Checking logs...${NC}"
    journalctl -u $APP_NAME -n 50 --no-pager
    echo ""
    echo "Please check the logs above for error details."
    exit 1
fi
