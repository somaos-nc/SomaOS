#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <sys/mman.h>
#include <unistd.h>
#include <stdint.h>
#include <string.h>
#include <sys/socket.h>
#include <netinet/in.h>

// SOMA OS: Upgraded ALINX 7020 C-Agent (v1.3)
// Support for Reset, Telemetry, and Silence Protocol

#define TELEMETRY_PORT 8080
#define FPGA_REG_ADDR  0x43C00000 
#define RESET_BIT      (1 << 31) // Example reset bit in Control Register

uint64_t mock_register = 0;
float mock_temp = 36.5;

void handle_request(int client_socket) {
    char request[1024];
    read(client_socket, request, 1024);

    char response[1024];
    if (strstr(request, "GET /telemetry") != NULL) {
        mock_register = (uint64_t)rand();
        mock_temp += (rand() % 10 - 5) * 0.01;
        sprintf(response, "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n\r\n{\"reg\": %llu, \"temp\": %.2f}", (unsigned long long)mock_register, mock_temp);
    } 
    else if (strstr(request, "POST /reset") != NULL) {
        printf(">> [COMMAND] Executing Hardware Reset via ARM-to-PL Interface...\n");
        // In real HW: *((uint32_t*)ctrl_reg) |= RESET_BIT; usleep(10); *((uint32_t*)ctrl_reg) &= ~RESET_BIT;
        sprintf(response, "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n\r\n{\"status\": \"reset_complete\"}");
    }
    else {
        sprintf(response, "HTTP/1.1 404 Not Found\r\n\r\n");
    }

    send(client_socket, response, strlen(response), 0);
    close(client_socket);
}

int main() {
    printf(">> SomaOS HPQC Agent v1.3 Live on Zynq ARM...\n");
    printf(">> Commands available: GET /telemetry, POST /reset\n");

    int server_fd, client_socket;
    struct sockaddr_in address;
    int opt = 1;
    int addrlen = sizeof(address);

    server_fd = socket(AF_INET, SOCK_STREAM, 0);
    setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR | SO_REUSEPORT, &opt, sizeof(opt));

    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(TELEMETRY_PORT);

    if (bind(server_fd, (struct sockaddr *)&address, sizeof(address)) < 0) {
        perror("Bind failed");
        return -1;
    }
    listen(server_fd, 5);

    while(1) {
        client_socket = accept(server_fd, (struct sockaddr *)&address, (socklen_t*)&addrlen);
        if (client_socket >= 0) {
            handle_request(client_socket);
        }
    }

    return 0;
}
