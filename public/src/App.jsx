import { createSignal } from "solid-js";

function App() {
  const [data, setData] = createSignal();
  const [isError, setIsError] = createSignal();

  function getData(ev) {
    setIsError(false);
    ev.preventDefault();
    const data = Object.fromEntries(new FormData(ev.currentTarget));

    const promise = fetch('http://localhost:8069/intoTripples', {
      method: 'POST',
      body: data.turtle,
    });
    
    promise.then(r => r.json()).then(d => {
        if ('error' in d) {
            setIsError(true);
            setData(d.error);
            return;
        }

        setData(JSON.stringify(d, null, 4));
    });
  }

  return (
    <div>
      <form 
        onSubmit={getData}
        class="p-4 flex flex-col gap-2 max-w-xl items-start"
      >
        <textarea 
          id="turtle-textarea"
          class="border-2 focus:border-green-500 shadow-md rounded-2xl overflow-hidden p-4 outline-none overflow-y-scroll"
          cols="60"
          rows="10"
          name="turtle"
          placeholder="@prefix : <http://tvoje.mama/> ."
        ></textarea>
        <button 
          type="submit"
          class="rounded-2xl border-2 bg-white border-green-500 p-4 shadow-md font-bold hover:brightness-95"
        >Konvertovat</button>
      </form>

      { isError() ? (
            <code class="text-red-500 font-bold px-4">{data()}</code>
      ) : (
            <pre class="px-4">{data()}</pre>
      )}
    </div>
  );
}

export default App;
