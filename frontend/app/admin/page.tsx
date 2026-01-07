import { getServerSession } from "next-auth";
import { redirect } from "next/navigation";

export default async function AdminPage() {
  const session = await getServerSession();

  if (!session || (session.user as any).role !== "admin") {
    redirect("/login");
  }

  return <h1 className="text-xl">Admin Panel</h1>;
}
