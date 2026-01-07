import NextAuth from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";

const handler = NextAuth({
  providers: [
    CredentialsProvider({
      name: "Credentials",
      credentials: {
        email: { label: "Email", type: "text" },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials) {
        // TEMP USERS (we will replace with DB later)
        if (
          credentials?.email === "student@uni.edu" &&
          credentials?.password === "student123"
        ) {
          return {
            id: "1",
            name: "Student User",
            email: "student@uni.edu",
            role: "student",
          };
        }

        if (
          credentials?.email === "admin@uni.edu" &&
          credentials?.password === "admin123"
        ) {
          return {
            id: "2",
            name: "Admin User",
            email: "admin@uni.edu",
            role: "admin",
          };
        }

        return null;
      },
    }),
  ],
  callbacks: {
    async jwt({ token, user }) {
      if (user) token.role = (user as any).role;
      return token;
    },
    async session({ session, token }) {
      if (session.user) {
        (session.user as any).role = token.role;
      }
      return session;
    },
  },
  pages: {
    signIn: "/login",
  },
});

export { handler as GET, handler as POST };
