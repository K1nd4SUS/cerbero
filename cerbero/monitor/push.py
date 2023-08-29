import time
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler
import os
import subprocess


class  MyHandler(FileSystemEventHandler):
   def  on_modified(self,  event):
      if event.src_path == "config.json":
         print(f'File modificato')
         subprocess.run(["rsync", "-v", "-r", "config.json", f"root@{os.getenv('VULNBOX_IP')}:/root/cerbero"])

if __name__ ==  "__main__":
   print("Creating cerbero folder...")
   # subprocess.run(["sshpass", "-p", f"{os.getenv('VULNBOX_PW')}", "ssh", "-o StrictHostKeyChecking=no", f"root@{os.getenv('VULNBOX_IP')}", "mkdir", "cerbero"])
   subprocess.run(["ssh", "-o StrictHostKeyChecking=no", f"root@{os.getenv('VULNBOX_IP')}", "mkdir", "cerbero"])
   print("Copy cerbero executable")
   # subprocess.run(["sshpass", "-p", f"{os.getenv('VULNBOX_PW')}", "rsync", "-v", "-r", "cerbero", f"root@{os.getenv('VULNBOX_IP')}:/root/cerbero"])
   subprocess.run(["rsync", "-v", "-r", "cerbero", f"root@{os.getenv('VULNBOX_IP')}:/root/cerbero"])
   print("Copy start_cerbero")
   # subprocess.run(["sshpass", "-p", f"{os.getenv('VULNBOX_PW')}", "rsync", "-v", "-r", "start_cerbero", f"root@{os.getenv('VULNBOX_IP')}:/root/cerbero"])
   subprocess.run(["rsync", "-v", "-r", "start_cerbero", f"root@{os.getenv('VULNBOX_IP')}:/root/cerbero"])
   print("Copy config.json")
   # subprocess.run(["sshpass", "-p", f"{os.getenv('VULNBOX_PW')}", "rsync", "-v", "-r", "config.json", f"root@{os.getenv('VULNBOX_IP')}:/root/cerbero"])
   subprocess.run(["rsync", "-v", "-r", "config.json", f"root@{os.getenv('VULNBOX_IP')}:/root/cerbero"])

   # ubprocess.run(["sshpass", "-p", f"{os.getenv('VULNBOX_PW')}", "ssh", "-o StrictHostKeyChecking=no", f"root@{os.getenv('VULNBOX_IP')}", "/root/cerbero/cerbero", "-t", "j", ">>", "log"]) 

   event_handler = MyHandler()
   observer = Observer()
   observer.schedule(event_handler,  path='config.json',  recursive=False)
   observer.start()

   try:
      while  True:
         time.sleep(1)
   except  KeyboardInterrupt:
      observer.stop()
      observer.join()

