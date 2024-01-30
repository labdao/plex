export default function Layout({ children, panel }: { children: React.ReactNode; panel: React.ReactNode }) {
  return (
    <div className="grid grid-flow-col auto-cols-fr">
      {children}
      {panel}
    </div>
  );
}
