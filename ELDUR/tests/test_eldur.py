import unittest
import sys
import os

# Add src to path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../src')))

from eldur import EldurCore

class TestEldurActiveDefense(unittest.TestCase):

    def setUp(self):
        self.initial_uid = "UID-VECT-ALPHA"
        self.eldur = EldurCore(initial_uid=self.initial_uid, threshold=0.2)

    def test_high_coherence_maintains_uid(self):
        """Test that normal, stable operations do not trigger a relocation."""
        uid_vector = [1.0, 1.0, 1.0]
        # Low entropy environment
        env_vector = [0.1, 0.2, 0.1]
        
        self.eldur.evaluate_state(uid_vector, env_vector)
        
        self.assertEqual(self.eldur.get_current_uid(), self.initial_uid)
        self.assertEqual(len(self.eldur.relocation_history), 1)

    def test_low_coherence_triggers_relocation(self):
        """Test that sub-quantum sabotage (high entropy/noise) triggers UID relocation."""
        uid_vector = [1.0, 1.0, 1.0]
        # High entropy environment (simulated sabotage / decoherence)
        env_vector = [10.0, 50.0, 20.0]
        
        self.eldur.evaluate_state(uid_vector, env_vector)
        
        # Verify UID has changed
        self.assertNotEqual(self.eldur.get_current_uid(), self.initial_uid)
        self.assertEqual(len(self.eldur.relocation_history), 2)
        self.assertTrue(self.eldur.get_current_uid().startswith("UID-VECT-"))

    def test_s_phi_calculation(self):
        """Test the S(Phi) mathematical calculation directly."""
        uid_vector = [1.0, 1.0]
        env_vector = [0.0, 0.0]
        
        # If lambda_const is 0.5, and env is 0:
        # S(phi) = (0.5/1 + 0.5/1) / 2 = 0.5
        s_phi = self.eldur.measure_scalar_coherence(uid_vector, env_vector)
        self.assertAlmostEqual(s_phi, 0.5)

if __name__ == '__main__':
    unittest.main()
