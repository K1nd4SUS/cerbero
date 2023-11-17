#!/bin/python3
import subprocess
import signal
import sys
import os
import time


cerbero_process: subprocess.Popen | None = None


def run_command(command: str, working_directory: str | None = None) -> subprocess.Popen:
    return subprocess.Popen(command, shell=True, cwd=working_directory)
    pass


def build_cerbero():
    # this command runs a specific task of docker compose.
    # the (~said~) aforementioned task is the last argument of the command
    run_command("docker compose run --build --rm --name builder cerbero2").wait()
    pass


def run_metrics_monitor():
    # this command runs a specific task of docker compose.
    # the (~said~) aforementioned task is the last argument of the command
    run_command(
        "docker compose run -d --build --service-ports --rm --name metrics_monitor metrics_monitor"
    ).wait()
    pass


def run_cerbero():
    global cerbero_process
    cerbero_process = run_command("./firewall2", "./cerbero2/bin")
    pass


def stop_metrics_monitor():
    run_command("docker stop metrics_monitor").wait()
    pass


def stop_cerbero():
    global cerbero_process
    cerbero_process.send_signal(signal.SIGINT)
    cerbero_process.wait()
    pass


def stopped_handler(sig, frame):
    print("Stopping metrics monitor...")
    stop_metrics_monitor()
    print("Stopped metrics monitor")

    print("Stopping metrics cerbero...")
    stop_cerbero()
    print("Stopped metrics cerbero")

    sys.exit(0)
    pass


if __name__ == "__main__":
    # root uid is 0
    if os.geteuid() != 0:
        print("This script must be run as root.")
        sys.exit(0)
        pass

    signal.signal(signal.SIGINT, stopped_handler)
    signal.signal(signal.SIGTERM, stopped_handler)

    print("Building cerbero...")
    build_cerbero()
    print("Built cerbero")

    print("Starting metrics monitor...")
    run_metrics_monitor()
    print("Started metrics monitor")

    print("Starting cerbero...")
    run_cerbero()
    print("Started cerbero")

    signal.pause()
    pass
