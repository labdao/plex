export default function Layout({ children, panel }: { children: React.ReactNode; panel: React.ReactNode }) {
  return (
    <div className="grid grid-flow-col auto-cols-fr lg:auto-cols-auto lg:grid-cols-[minmax(auto,_1fr)]">
      {children}
      {panel}
    </div>
  );
}
