import unittest
from fastapi.testclient import TestClient
import sys
import os

# Mock Vertex AI before importing router
from unittest.mock import MagicMock, patch
sys.modules['vertexai'] = MagicMock()
sys.modules['vertexai.generative_models'] = MagicMock()

# Add src to path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../src')))

from router import app

class TestSomaAIRouter(unittest.TestCase):
    def setUp(self):
        self.client = TestClient(app)

    def test_root_endpoint(self):
        """Test the AI Router's basic health check."""
        response = self.client.get("/")
        self.assertEqual(response.status_code, 200)
        self.assertIn("SomaAI Router is Online", response.json()["status"])

    @patch("router.model.generate_content_async")
    def test_api_telemetry_analysis(self, mock_gen):
        """Test that the telemetry analysis endpoint routes to Gemini."""
        # Mock the async response from Vertex AI
        mock_resp = MagicMock()
        mock_resp.text = "Mocked AI Insight: System is stable."
        
        # We need to mock the awaitable
        async def mock_async_gen(*args, **kwargs):
            return mock_resp
        
        mock_gen.side_effect = mock_async_gen
        
        payload = {"logs": "Sample hardware logs 123"}
        response = self.client.post("/api/ai/telemetry", json=payload)
        
        self.assertEqual(response.status_code, 200)
        self.assertIn("insight", response.json())
        self.assertEqual(response.json()["insight"], "Mocked AI Insight: System is stable.")

if __name__ == "__main__":
    unittest.main()
