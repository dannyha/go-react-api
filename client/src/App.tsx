import "./App.css";
import CreateTransation from "./containers/CreateTransation";

const App: React.FC = (): React.ReactElement => {
  return (
    <div className="app">
      <CreateTransation />
    </div>
  );
};

export default App;
