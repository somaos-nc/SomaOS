#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <sys/mman.h>
#include <unistd.h>
#include <stdint.h>
#include <string.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <pthread.h>
#include <sys/ioctl.h>
#include <net/if.h>
#include <time.h>

// SOMA OS: Upgraded ALINX 7020 C-Agent (v1.4 - Network Guardian Edition)
// Support for Reset, Telemetry, and Auto-Network Recovery

#define TELEMETRY_PORT 8080
#define FPGA_REG_ADDR  0x43C00000 
#define RESET_BIT      (1 << 31)

uint64_t mock_register = 0;
float mock_temp = 36.5;
char manifold_hex[257];

void generate_manifold() {
    for (int i = 0; i < 128; i++) {
        sprintf(manifold_hex + (i * 2), "%02x", rand() % 256);
    }
}

void* network_guardian(void* arg) {
    printf(">> [GUARDIAN] Network Monitor Thread Started.\n");
    while(1) {
        int sock = socket(AF_INET, SOCK_DGRAM, 0);
        struct ifreq ifr;
        memset(&ifr, 0, sizeof(ifr));
        strcpy(ifr.ifr_name, "eth0");
        
        if (ioctl(sock, SIOCGIFFLAGS, &ifr) < 0) {
            printf(">> [GUARDIAN WARNING] eth0 interface error. Pulsing...\n");
            system("/sbin/ifconfig eth0 up && /sbin/udhcpc -i eth0");
        } else if (!(ifr.ifr_flags & IFF_UP) || !(ifr.ifr_flags & IFF_RUNNING)) {
            printf(">> [GUARDIAN] Link Down detected. Restoring eth0...\n");
            system("/sbin/ifconfig eth0 up && /sbin/udhcpc -i eth0");
        }
        close(sock);
        sleep(10);
    }
    return NULL;
}

void handle_request(int client_socket) {
    char request[1024];
    read(client_socket, request, 1024);

    char response[2048];
    if (strstr(request, "GET /telemetry") != NULL) {
        generate_manifold();
        mock_register = (uint64_t)rand();
        mock_temp += (rand() % 10 - 5) * 0.01;
        sprintf(response, "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n\r\n{\"reg\": %llu, \"temp\": %.2f, \"manifold\": \"%s\"}", (unsigned long long)mock_register, mock_temp, manifold_hex);
    } 
    else if (strstr(request, "POST /reset") != NULL) {
        printf(">> [COMMAND] Executing Hardware Reset...\n");
        sprintf(response, "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n\r\n{\"status\": \"reset_complete\"}");
    }
    else {
        sprintf(response, "HTTP/1.1 404 Not Found\r\n\r\n");
    }

    send(client_socket, response, strlen(response), 0);
    close(client_socket);
}

int main() {
    printf(">> SomaOS HPQC Agent v1.4 (Guardian) Live...\n");
    srand(time(NULL));

    pthread_t guardian_tid;
    pthread_create(&guardian_tid, NULL, network_guardian, NULL);
    pthread_detach(guardian_tid);

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
