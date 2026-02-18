export default function Home() {
	return (
		<div className="flex flex-col gap-8 my-8">
			<p className="border-b border-b-black pb-8">
				Welcome to the <strong>Monte-Markov Weather Simulation</strong>, where we aim to estimate the maximum and minimum temperatures and the likelihood of rain for the next 3 days, using only Markov Chains and Monte Carlo.
			</p>
			<form className="flex flex-col gap-2">
				<label className="font-bold text-lg">Type your current location</label>
				<input className="border-black border rounded-md px-4 py-2" />
			</form>
		</div>
	);
}
