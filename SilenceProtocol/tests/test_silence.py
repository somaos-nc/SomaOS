import unittest
import sys
import os

# Add src to path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../src')))

from silence import SilenceProtocol

class TestSilenceProtocol(unittest.TestCase):
    def setUp(self):
        # Base unit of 100ms
        self.protocol = SilenceProtocol(base_time_unit=0.1)

    def test_encode_decode_symmetric(self):
        """Test that data encoded into silence intervals can be perfectly decoded."""
        original_data = [3, 0, 1, 2, 1]
        
        # Encode data to intervals
        intervals = self.protocol.encode(original_data)
        
        # Simulate network arrival timestamps
        timestamps = [0.0]
        current_time = 0.0
        for interval in intervals:
            current_time += interval
            timestamps.append(current_time)

        # Decode
        decoded_data = self.protocol.decode(timestamps)
        
        self.assertEqual(original_data, decoded_data)

    def test_parity_generation(self):
        """Test the mathematically specified quaternary error correction: P = (sum S_i) mod 4"""
        data = [2, 1, 2] # Sum = 5, 5 mod 4 = 1
        intervals = self.protocol.encode(data)
        
        # The last interval should be the parity ping (1 * 0.1s)
        self.assertAlmostEqual(intervals[-1], 0.1)

    def test_parity_error_detection(self):
        """Test that tampered intervals trigger a parity exception."""
        timestamps = [0.0, 0.2, 0.3, 0.5, 0.8] # Symbols: 2, 1, 2, Parity=3
        # Valid parity would be (2+1+2) mod 4 = 1. We supplied 3 (0.8 - 0.5 = 0.3s).
        
        with self.assertRaises(ValueError) as context:
            self.protocol.decode(timestamps)
        
        self.assertTrue("Parity Error" in str(context.exception))

    def test_jitter_tolerance(self):
        """Test the system's ability to decode signals with thermal/network jitter."""
        original_data = [1, 3]
        
        # Simulate timestamps with slight delays
        # Perfect: 0.0 -> 0.1 (1) -> 0.4 (3) -> 0.4 (Parity 0)
        timestamps = [0.0, 0.103, 0.395, 0.402]
        
        decoded_data = self.protocol.decode(timestamps)
        self.assertEqual(original_data, decoded_data)

if __name__ == '__main__':
    unittest.main()
