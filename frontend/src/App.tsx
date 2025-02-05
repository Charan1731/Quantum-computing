import React, { useState } from 'react';
import { KeySquare, Shield, CheckCircle2, AlertCircle, Edit3 } from 'lucide-react';

interface VerificationResult {
  valid: boolean;
}

function App() {
  const [publicKey, setPublicKey] = useState('');
  const [message, setMessage] = useState('');
  const [signature, setSignature] = useState('');
  const [verificationResult, setVerificationResult] = useState<VerificationResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [isEditingPublicKey, setIsEditingPublicKey] = useState(false);
  const [isEditingSignature, setIsEditingSignature] = useState(false);

  const generateKeys = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await fetch('http://localhost:8080/generate-key');
      const data = await response.json();
      setPublicKey(data.publicKey);
    } catch (err) {
      setError('Failed to generate keys. Is the backend running?');
    } finally {
      setLoading(false);
    }
  };

  const signMessage = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await fetch('http://localhost:8080/sign', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ message }),
      });
      const data = await response.json();
      setSignature(data.signature);
    } catch (err) {
      setError('Failed to sign message. Make sure to generate keys first.');
    } finally {
      setLoading(false);
    }
  };

  const verifySignature = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await fetch('http://localhost:8080/verify', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          message,
          signature,
          publicKey,
        }),
      });
      const data = await response.json();
      setVerificationResult(data);
    } catch (err) {
      setError('Failed to verify signature.');
    } finally {
      setLoading(false);
    }
  };

  const renderEditableField = (
    value: string,
    isEditing: boolean,
    setIsEditing: (value: boolean) => void,
    onChange: (value: string) => void,
    label: string
  ) => (
    <div className="relative">
      <p className="text-sm text-gray-400 mb-2 flex items-center justify-between">
        <span>{label}</span>
        <button
          onClick={() => setIsEditing(!isEditing)}
          className="text-blue-400 hover:text-blue-300 transition-colors flex items-center gap-1 text-sm"
        >
          <Edit3 className="w-4 h-4" />
          {isEditing ? 'Done' : 'Edit'}
        </button>
      </p>
      {isEditing ? (
        <textarea
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2 text-white font-mono text-sm min-h-[100px]"
        />
      ) : (
        <p className="font-mono text-sm bg-gray-900 p-3 rounded break-all">
          {value || 'No value set'}
        </p>
      )}
    </div>
  );

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 to-gray-800 text-white">
      <div className="container mx-auto px-4 py-12 max-w-4xl">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold mb-4 flex items-center justify-center gap-3">
            <Shield className="w-10 h-10" />
            Quantum-Safe Blockchain Signatures
          </h1>
          <p className="text-gray-400">
            Secure your blockchain transactions with quantum-resistant cryptography
          </p>
        </div>

        {error && (
          <div className="bg-red-900/50 border border-red-500 rounded-lg p-4 mb-6 flex items-center gap-3">
            <AlertCircle className="w-5 h-5 text-red-500" />
            <p className="text-red-200">{error}</p>
          </div>
        )}

        <div className="grid gap-6">
          <div className="bg-gray-800/50 p-6 rounded-lg border border-gray-700">
            <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
              <KeySquare className="w-5 h-5" /> Key Generation
            </h2>
            <button
              onClick={generateKeys}
              disabled={loading}
              className="bg-blue-600 hover:bg-blue-700 px-4 py-2 rounded-lg transition-colors disabled:opacity-50"
            >
              Generate New Key Pair
            </button>
            {renderEditableField(
              publicKey,
              isEditingPublicKey,
              setIsEditingPublicKey,
              setPublicKey,
              'Public Key'
            )}
          </div>

          <div className="bg-gray-800/50 p-6 rounded-lg border border-gray-700">
            <h2 className="text-xl font-semibold mb-4">Message Signing</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm text-gray-400 mb-2">
                  Message to Sign
                </label>
                <input
                  type="text"
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                  className="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2 text-white"
                  placeholder="Enter your message"
                />
              </div>
              <button
                onClick={signMessage}
                disabled={loading || !message}
                className="bg-green-600 hover:bg-green-700 px-4 py-2 rounded-lg transition-colors disabled:opacity-50"
              >
                Sign Message
              </button>
              {renderEditableField(
                signature,
                isEditingSignature,
                setIsEditingSignature,
                setSignature,
                'Signature'
              )}
            </div>
          </div>

          <div className="bg-gray-800/50 p-6 rounded-lg border border-gray-700">
            <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
              <CheckCircle2 className="w-5 h-5" /> Signature Verification
            </h2>
            <button
              onClick={verifySignature}
              disabled={loading || !message || !signature || !publicKey}
              className="bg-purple-600 hover:bg-purple-700 px-4 py-2 rounded-lg transition-colors disabled:opacity-50"
            >
              Verify Signature
            </button>
            {verificationResult !== null && (
              <div className={`mt-4 p-4 rounded-lg ${
                verificationResult.valid 
                  ? 'bg-green-900/30 border border-green-600' 
                  : 'bg-red-900/30 border border-red-600'
              }`}>
                <p className="flex items-center gap-2">
                  {verificationResult.valid ? (
                    <>
                      <CheckCircle2 className="w-5 h-5 text-green-500" />
                      <span className="text-green-400">Signature is valid!</span>
                    </>
                  ) : (
                    <>
                      <AlertCircle className="w-5 h-5 text-red-500" />
                      <span className="text-red-400">Invalid signature!</span>
                    </>
                  )}
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;