import math
from typing import List

class SilenceProtocol:
    """
    Implementation of The Silence Protocol: A quaternary exolinguistic network topology.
    Unlike standard protocols (TCP/IP) which encode data into the electromagnetic
    carrier wave itself, the Silence Protocol uses the carrier wave purely as a 
    delimiter (a "ping"). The actual payload is encoded in the duration of SILENCE
    between the pings.
    """

    def __init__(self, base_time_unit: float = 0.1):
        """
        Initialize the protocol.
        :param base_time_unit: The atomic unit of time (Delta t) in seconds.
        """
        self.base_time_unit = base_time_unit

    def encode(self, data: List[int]) -> List[float]:
        """
        Encodes a list of Base-4 integers into a series of time intervals.
        Each interval represents the silence before the next ping.
        
        :param data: List of integers [0, 1, 2, 3].
        :return: List of time intervals (silences) in seconds.
        """
        intervals = []
        for symbol in data:
            if not (0 <= symbol <= 3):
                raise ValueError(f"Invalid symbol {symbol}. Protocol is strictly Quaternary (Base-4).")
            
            # The symbol (0-3) determines the multiplier for the base time unit.
            # E.g., symbol '2' -> wait 2 * base_time_unit seconds.
            interval = symbol * self.base_time_unit
            intervals.append(interval)

        # Append Quaternary Error Correction Ping
        # P = (\sum S_i) mod 4
        parity = sum(data) % 4
        intervals.append(parity * self.base_time_unit)

        return intervals

    def decode(self, timestamps: List[float], tolerance: float = 0.05) -> List[int]:
        """
        Decodes a stream of ping timestamps back into Base-4 data.
        
        :param timestamps: Absolute times (in seconds) when pings were received.
        :param tolerance: Acceptable variance in network jitter.
        :return: Decoded Quaternary data payload.
        """
        if len(timestamps) < 2:
            return []

        # 1. Calculate the intervals (silence durations) between pings
        intervals = []
        for i in range(1, len(timestamps)):
            intervals.append(timestamps[i] - timestamps[i-1])

        # 2. Re-derive the atomic unit of time Delta t dynamically (syncing to sender)
        # Using a simplified GCD heuristic for floating point network timing
        # Delta t = GCD(T_1, T_2, ..., T_n)
        # For this prototype, we'll assume the base_time_unit is roughly known, 
        # but in a true implementation, it is derived purely from the stream geometry.
        inferred_delta_t = self.base_time_unit 

        # 3. Decode intervals into symbols
        data = []
        for interval in intervals:
            # Round to the nearest integer multiple of Delta t
            symbol = round(interval / inferred_delta_t)
            
            # Constraint check (Signal Integrity)
            if abs((symbol * inferred_delta_t) - interval) > tolerance:
                raise ValueError("Signal integrity compromised: Jitter exceeds tolerance.")
            
            data.append(int(symbol))

        if not data:
            return []

        # 4. Validate Parity
        payload = data[:-1]
        parity_ping = data[-1]
        
        expected_parity = sum(payload) % 4
        if parity_ping != expected_parity:
            raise ValueError(f"Parity Error: Expected {expected_parity}, Got {parity_ping}")

        return payload
