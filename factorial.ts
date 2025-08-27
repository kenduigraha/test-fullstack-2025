function factorial(n: number): number {
  if (n < 0) return 1;
  let result = 1;
  for (let i = 2; i <= n; i++) {
    result *= i;
  }
  return result;
}

function f(n: number): number {
  if (n < 0) return 1;

  const numerator = factorial(n);
  const denominator = 2 ** n;

  const division = numerator / denominator;
  const remainder = numerator % denominator;

  console.log('numerator: ', numerator)
  console.log('denominator: ', denominator)
  console.log('division: ', division)

  // Jika ada sisa, berarti perlu pembulatan ke atas
  return remainder === 0 ? division : Math.ceil(division);
}

// contoh
console.log(f(0)); // Output: 1 
console.log(f(5)); // Output: 4
console.log(f(10)); // Output: 3544
