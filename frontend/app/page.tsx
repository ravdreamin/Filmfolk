export default function Home() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-8">
      <main className="flex flex-col items-center gap-8 max-w-4xl">
        <h1 className="text-5xl font-bold text-center">
          Welcome to FilmFolk
        </h1>
        <p className="text-xl text-center text-gray-600 dark:text-gray-400">
          Discover, review, and share your favorite movies with the community
        </p>
        <div className="flex gap-4">
          <a
            href="/movies"
            className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
          >
            Browse Movies
          </a>
          <a
            href="/login"
            className="px-6 py-3 border border-gray-300 dark:border-gray-700 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition"
          >
            Sign In
          </a>
        </div>
      </main>
    </div>
  );
}
