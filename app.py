def calc(n):
  if n <= 1:
    return n
  else:
    return calc(n-1) + calc(n-2)
    
print(calc(6))