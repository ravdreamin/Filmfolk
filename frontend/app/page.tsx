'use client'

import { motion } from 'framer-motion'
import Link from 'next/link'
import { Film, Star, Users, MessageCircle, TrendingUp, Sparkles } from 'lucide-react'
import Button from '@/components/ui/Button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card'

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
    },
  },
}

const itemVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.5,
    },
  },
}

export default function Home() {
  return (
    <div className="flex flex-col min-h-screen">
      {/* Hero Section */}
      <section className="relative overflow-hidden py-20 sm:py-32">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <motion.div
            className="text-center"
            initial={{ opacity: 0, y: 50 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
          >
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ delay: 0.2, type: 'spring', stiffness: 200 }}
              className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 border border-primary/20 mb-8"
            >
              <Sparkles className="h-4 w-4 text-primary" />
              <span className="text-sm font-medium text-primary">Your Movie Community Awaits</span>
            </motion.div>

            <h1 className="text-4xl sm:text-6xl lg:text-7xl font-bold tracking-tight mb-6">
              <span className="block text-foreground">Discover Movies.</span>
              <span className="block bg-gradient-to-r from-primary via-purple-400 to-pink-400 bg-clip-text text-transparent">
                Share Reviews.
              </span>
              <span className="block text-foreground">Connect with Cinephiles.</span>
            </h1>

            <p className="text-xl text-muted-foreground max-w-2xl mx-auto mb-8">
              Join FilmFolk, the ultimate social platform for movie lovers. Discover new films,
              write thoughtful reviews, and connect with people who share your passion for cinema.
            </p>

            <motion.div
              className="flex flex-col sm:flex-row gap-4 justify-center"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.6 }}
            >
              <Link href="/register">
                <Button size="lg" className="w-full sm:w-auto">
                  Get Started Free
                </Button>
              </Link>
              <Link href="/movies">
                <Button size="lg" variant="outline" className="w-full sm:w-auto">
                  Explore Movies
                </Button>
              </Link>
            </motion.div>

            {/* Stats */}
            <motion.div
              className="grid grid-cols-3 gap-8 max-w-2xl mx-auto mt-16"
              variants={containerVariants}
              initial="hidden"
              animate="visible"
            >
              <motion.div variants={itemVariants}>
                <div className="text-3xl font-bold text-primary">10K+</div>
                <div className="text-sm text-muted-foreground">Movies</div>
              </motion.div>
              <motion.div variants={itemVariants}>
                <div className="text-3xl font-bold text-primary">50K+</div>
                <div className="text-sm text-muted-foreground">Reviews</div>
              </motion.div>
              <motion.div variants={itemVariants}>
                <div className="text-3xl font-bold text-primary">5K+</div>
                <div className="text-sm text-muted-foreground">Users</div>
              </motion.div>
            </motion.div>
          </motion.div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 bg-muted/50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <motion.div
            className="text-center mb-16"
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
          >
            <h2 className="text-3xl sm:text-4xl font-bold mb-4">
              Everything You Need to Love Movies
            </h2>
            <p className="text-muted-foreground text-lg max-w-2xl mx-auto">
              FilmFolk brings together all the tools you need to track, review, and discuss your favorite films.
            </p>
          </motion.div>

          <motion.div
            className="grid md:grid-cols-2 lg:grid-cols-3 gap-6"
            variants={containerVariants}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
          >
            <FeatureCard
              icon={<Film className="h-8 w-8" />}
              title="Extensive Movie Database"
              description="Access a comprehensive database of movies with detailed information, cast details, and ratings powered by TMDB."
            />
            <FeatureCard
              icon={<Star className="h-8 w-8" />}
              title="Write Reviews"
              description="Share your thoughts with detailed reviews and ratings. Help others discover great films."
            />
            <FeatureCard
              icon={<MessageCircle className="h-8 w-8" />}
              title="Threaded Discussions"
              description="Engage in meaningful conversations with threaded comments on reviews and movies."
            />
            <FeatureCard
              icon={<Users className="h-8 w-8" />}
              title="Build Your Network"
              description="Connect with fellow movie enthusiasts and discover new films based on shared tastes."
            />
            <FeatureCard
              icon={<TrendingUp className="h-8 w-8" />}
              title="Track Your Journey"
              description="Keep lists of movies you've watched, want to watch, and dropped. See your viewing stats over time."
            />
            <FeatureCard
              icon={<Sparkles className="h-8 w-8" />}
              title="Smart Recommendations"
              description="Get personalized movie recommendations based on your viewing history and preferences."
            />
          </motion.div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <motion.div
            className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-primary to-purple-600 p-12 text-center"
            initial={{ opacity: 0, scale: 0.9 }}
            whileInView={{ opacity: 1, scale: 1 }}
            viewport={{ once: true }}
          >
            <div className="relative z-10">
              <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
                Ready to Join the Community?
              </h2>
              <p className="text-white/90 text-lg mb-8 max-w-2xl mx-auto">
                Start your journey with FilmFolk today. It's free, it's fun, and it's filled with people who love movies as much as you do.
              </p>
              <Link href="/register">
                <Button size="lg" variant="secondary">
                  Create Your Free Account
                </Button>
              </Link>
            </div>
          </motion.div>
        </div>
      </section>
    </div>
  )
}

function FeatureCard({ icon, title, description }: { icon: React.ReactNode; title: string; description: string }) {
  return (
    <motion.div variants={itemVariants}>
      <Card className="h-full">
        <CardHeader>
          <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center text-primary mb-4">
            {icon}
          </div>
          <CardTitle className="text-xl">{title}</CardTitle>
        </CardHeader>
        <CardContent>
          <CardDescription className="text-base">{description}</CardDescription>
        </CardContent>
      </Card>
    </motion.div>
  )
}
