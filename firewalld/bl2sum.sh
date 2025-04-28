#!/bin/bash

# Log function to print messages with timestamp
log_info() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] $1"
}

log_error() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] $1"
}

# File to check
FILE="firewalld2.go"

# Check if the file exists
if [[ -f "$FILE" ]]; then
    log_info "File '$FILE' found, calculating BLAKE2 checksum..."
    
    # Calculate the BLAKE2 sum and log it
    checksum=$(b2sum "$FILE")
    
    if [[ $? -eq 0 ]]; then
        log_info "BLAKE2 checksum for '$FILE': $checksum"
    else
        log_error "Error calculating BLAKE2 checksum for '$FILE'."
    fi
else
    log_error "File '$FILE' not found."
    exit 1
fi
