import { PORTS } from '@twir/grpc/constants/constants';
import * as Eval from '@twir/grpc/generated/eval/eval';
import _ from 'lodash';
import { createServer } from 'nice-grpc';
import { VM } from 'vm2';

const vm = new VM({
	sandbox: {
		fetch,
		URLSearchParams,
		_: _,
	},
	timeout: 1000,
	wasm: false,
	eval: false,
});

const evalService: Eval.EvalServiceImplementation = {
	async process(request: Eval.Evaluate): Promise<Eval.DeepPartial<Eval.EvaluateResult>> {
		let resultOfExecution: any;
		try {
			const toEval = `(async function () { ${request.script} })()`.split(';\n').join(';');
			resultOfExecution = await vm.run(toEval);
		} catch (error) {
			console.error(error);
			resultOfExecution = (error as any).message ?? 'unexpected error';
		}

		return {
			result: String(resultOfExecution).slice(0, 5000),
		};
	},
};

const server = createServer({
	'grpc.keepalive_time_ms': 1 * 60 * 1000,
});

server.add(Eval.EvalDefinition, evalService);

await server.listen(`0.0.0.0:${PORTS.EVAL_SERVER_PORT}`);
console.log('Eval microservice started');
