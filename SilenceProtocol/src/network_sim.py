import socket
import time
import threading
from silence import SilenceProtocol

# Network Config
UDP_IP = "127.0.0.1"
UDP_PORT = 5005
PING_PAYLOAD = b""  # The Silence Protocol sends exactly 0 bytes.

class SilenceSender:
    def __init__(self, base_time_unit=0.2):
        self.protocol = SilenceProtocol(base_time_unit)
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    def transmit(self, data: list):
        print(f"[SENDER] Preparing to transmit Quaternary Data: {data}")
        intervals = self.protocol.encode(data)
        
        print("[SENDER] Initiating transmission sequence...")
        
        # Send initial sync ping
        self.sock.sendto(PING_PAYLOAD, (UDP_IP, UDP_PORT))
        print(f"[SENDER] Sync ping sent.")

        for i, wait_time in enumerate(intervals):
            # The payload is the silence
            print(f"[SENDER] SILENCE for {wait_time:.2f}s (Symbol {wait_time / self.protocol.base_time_unit:.0f})")
            time.sleep(wait_time)
            
            # Send delimiter ping
            self.sock.sendto(PING_PAYLOAD, (UDP_IP, UDP_PORT))
            print(f"[SENDER] Delimiter ping {i+1} sent.")

        print("[SENDER] Transmission Complete.")

class SilenceReceiver:
    def __init__(self, base_time_unit=0.2):
        self.protocol = SilenceProtocol(base_time_unit)
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        self.sock.bind((UDP_IP, UDP_PORT))
        self.sock.settimeout(5.0) # Stop listening after 5 seconds of silence

    def listen(self):
        print(f"[RECEIVER] Listening on {UDP_IP}:{UDP_PORT}...")
        timestamps = []
        
        try:
            while True:
                data, addr = self.sock.recvfrom(1024)
                # Record the exact time the ping was received
                arrival_time = time.time()
                timestamps.append(arrival_time)
                print(f"[RECEIVER] Ping received at {arrival_time:.4f}")
        except socket.timeout:
            print("[RECEIVER] Stream ended (Timeout).")
            
        if len(timestamps) > 1:
            print("[RECEIVER] Decoding geometric silence...")
            try:
                decoded_data = self.protocol.decode(timestamps, tolerance=0.08)
                print(f"[RECEIVER] SUCCESS: Decoded Data -> {decoded_data}")
            except Exception as e:
                print(f"[RECEIVER] DECODE ERROR: {e}")
        else:
            print("[RECEIVER] Not enough pings to form a geometry.")

if __name__ == "__main__":
    # A simple demonstration sending Quaternary data
    # Data represents a simple state [2, 0, 1, 3]
    test_data = [2, 0, 1, 3]
    
    receiver = SilenceReceiver(base_time_unit=0.2)
    sender = SilenceSender(base_time_unit=0.2)

    # Start receiver in a background thread
    listener_thread = threading.Thread(target=receiver.listen)
    listener_thread.start()

    # Give the receiver a moment to bind to the port
    time.sleep(0.5)

    # Send the data
    sender.transmit(test_data)

    # Wait for the receiver to finish (timeout)
    listener_thread.join()
