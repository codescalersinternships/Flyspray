const R_PASSWORD = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
export type ValidationResult = {
  isValid: boolean;
  errorMessage: string;
};

export function validateEmail(value: string): ValidationResult {
  if (!value || value.trim() === "") {
    return {
      isValid: false,
      errorMessage: "Email is required",
    };
  }

  if (!R_PASSWORD.test(value)) {
    return {
      isValid: false,
      errorMessage: "Invalid email format",
    };
  }

  return {
    isValid: true,
    errorMessage: "",
  };
}

//isEightCharsLong checks if password is at least eight characters
function isEightCharsLong(value: string): boolean {
  if (value.length >= 8) return true;
  return false;
}

//containsOneLetter checks if password contains at least one letter
function containsOneLetter(value: string): boolean {
  const regexPattern = /[A-Za-z]/;
  return regexPattern.test(value);
}

//containsOneNumber checks if password contains at least one number
function containsOneNumber(value: string): boolean {
  const regexPattern = /\d/;
  return regexPattern.test(value);
}

//containsOneSymbol checks if password contains at least one special character
function containsOneSymbol(value: string): boolean {
  const regexPattern = /[_@$!%*#?&]/;
  return regexPattern.test(value);
}
export function validatePassword(value: string): ValidationResult {
  if (!value || value.trim() === "") {
    return {
      isValid: false,
      errorMessage: "Password is required",
    };
  }
  if (!isEightCharsLong(value)) {
    return {
      isValid: false,
      errorMessage: "Password must be at least eight characters long",
    };
  }
  if (!containsOneLetter(value)) {
    return {
      isValid: false,
      errorMessage: "Password must contain at least one letter",
    };
  }
  if (!containsOneNumber(value)) {
    return {
      isValid: false,
      errorMessage: "Password must contain at least one number",
    };
  }
  if (!containsOneSymbol(value)) {
    return {
      isValid: false,
      errorMessage: "Password must contain at least one special character",
    };
  }
  return {
    isValid: true,
    errorMessage: "",
  };
}

export function validateUsername(value: string): ValidationResult {
  if (!value || value.trim() === "") {
    return {
      isValid: false,
      errorMessage: "Username is required",
    };
  }
  if (!(value.length >= 4 && value.length <= 20)) {
    return {
      isValid: false,
      errorMessage: "Username should be four to twenty characters long",
    };
  }
  //to be implemented: add check that username does not already exist
  return {
    isValid: true,
    errorMessage: "",
  };
}
