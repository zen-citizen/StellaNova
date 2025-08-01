import clsx from "clsx";
import { useState } from "react";

function App() {
  const [count, setCount] = useState(0);

  return (
    <div className="container px-20 py-24">
      <h1
        className={clsx([
          "bg-gradient-to-r from-indigo-500 to-amber-600",
          "text-6xl font-bold tracking-wider",
          "uppercase",
          "w-max bg-clip-text text-transparent",
        ])}
      >
        Civic Compass
      </h1>
    </div>
  );
}

export default App;
