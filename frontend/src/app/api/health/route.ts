import { NextResponse } from 'next/server';

export async function GET() {
  try {
    const backendUrl = process.env.BACKEND_INTERNAL_URL;
    if (!backendUrl) {
      return new NextResponse('Error: BACKEND_INTERNAL_URL is not defined in Next.js API Route.', { status: 500 });
    }

    const response = await fetch(`${backendUrl}/api/health`);

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP Error Status: ${response.status}, Body: ${errorText}`);
    }

    const text = await response.text();
    return new NextResponse(text);
  } catch (error: any) {
    console.error('Error in Next.js API Route:', error.message || error);
    return new NextResponse(`Internal Server Error in Next.js API Route: ${error.message || 'Unknown error'}`, { status: 500 });
  }
}