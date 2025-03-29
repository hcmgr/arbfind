import { SignedIn, SignedOut, SignInButton, UserButton } from "@clerk/clerk-react";
import { useEffect, useState } from "react";

function ClerkStuff() {
  return (
    <header>
      <SignedOut>
        <SignInButton />
      </SignedOut>
      <SignedIn>
        <UserButton />
      </SignedIn>
    </header>
  );
}

type Outcome = {
    name:           string;
    price:          number;
}

type Arb = {
    sportkey:       string;
    sportitle:      string;
    commencetime:   string;
    hometeam:       string;
    awayteam:       string;
    outcomes:       Outcome[];
    r:              number;
};
  
type ArbBoxProps = {
    arb: Arb;
    onClick: (arb: Arb) => void;
};

function ArbBox({ arb, onClick }: ArbBoxProps) {
    const profitPercentage = Math.round((1 - arb.r) * 10000) / 100;
  
    return (
      <div className="arb-box" onClick={() => onClick(arb)}>
        <div><strong>{arb.sportkey}</strong></div>
        <div>{arb.hometeam} vs {arb.awayteam}</div>
        <div style={{ color: 'green' }}>+{profitPercentage}%</div>
      </div>
    );
  }

type DialogProps = {
    arb: Arb | null;
    onClose: () => void;
}

/**
 * Dialog component
 */
function Dialog({ arb, onClose }: DialogProps) {
  if (!arb) return null;
  const profitPercentage = Math.round((1 - arb.r) * 10000) / 100;

  return (
    <div className="dialog-backdrop" onClick={onClose}>
      <div className="dialog" onClick={e => e.stopPropagation()}>
        <div className="close-btn" onClick={onClose}>×</div>
        <h2>{arb.sportitle}</h2>
        <p>{arb.hometeam} vs {arb.awayteam}</p>
        <p style={{ color: 'green' }}>+{profitPercentage}%</p>
      </div>
    </div>
  );
}

/**
 * Core App component
 */
function App() {
  const [arbs, setArbs] = useState([]);
  const [selected, setSelected] = useState<Arb | null>(null);


  useEffect(() => {
    fetch('http://localhost:10000/getArbs')
      .then(res => res.json())
      .then(setArbs)
      .catch(console.error);
  }, []);

  return (
    <div className="container">
      <ClerkStuff />
      {arbs.map((arb, i) => (
        <ArbBox key={i} arb={arb} onClick={setSelected} />
      ))}
      <Dialog arb={selected} onClose={() => setSelected(null)} />
    </div>
  );
}

export default App
