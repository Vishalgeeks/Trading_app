import { Link } from 'react-router-dom';

export default function Home() {
  return (
    <div className="bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <div className="text-center">
          <h1 className="text-4xl md:text-6xl font-bold text-gray-900 mb-6">
            Beautiful Henna Art for Your <span className="text-rose-500">Special Day</span>
          </h1>
          <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">
            Book stunning mehndi designs online. Easy scheduling, beautiful designs, and unforgettable moments.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link to="/designs" className="btn-primary text-lg px-8 py-3">
              Browse Designs
            </Link>
            <Link to="/bookings" className="btn-secondary text-lg px-8 py-3">
              My Bookings
            </Link>
          </div>
        </div>

        <div className="mt-20 grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="card text-center">
            <div className="text-4xl mb-4">🎨</div>
            <h3 className="text-xl font-semibold mb-2">Beautiful Designs</h3>
            <p className="text-gray-600">Choose from our curated collection of traditional and modern mehndi patterns.</p>
          </div>
          <div className="card text-center">
            <div className="text-4xl mb-4">📅</div>
            <h3 className="text-xl font-semibold mb-2">Easy Booking</h3>
            <p className="text-gray-600">Select your preferred date and time slot. We'll confirm your appointment instantly.</p>
          </div>
          <div className="card text-center">
            <div className="text-4xl mb-4">✨</div>
            <h3 className="text-xl font-semibold mb-2">Professional Service</h3>
            <p className="text-gray-600">Experience premium henna application with high-quality, natural ingredients.</p>
          </div>
        </div>
      </div>
    </div>
  );
}
