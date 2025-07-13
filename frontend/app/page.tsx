
'use client';

import { useEffect, useState } from 'react';

export default function HomePage() {
  const [message, setMessage] = useState('');

  useEffect(() => {
    fetch('/api')
      .then((res) => res.text())
      .then((text) => setMessage(text));
  }, []);

  return (
    <div>
      <h1>Next.js Frontend</h1>
      <p>Database Ping Result: {message}</p>
    </div>
  );
}
