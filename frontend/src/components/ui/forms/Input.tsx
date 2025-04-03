import React from "react";

type InputProps = {
    label: string;
    type?: string;
    value: string;
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    required?: boolean;
    className?: string;
    placeholder?: string;
    error?: string;
};

export const Input = ({
    label,
    type = 'text',
    value,
    onChange,
    required = false,
    className = '',
    placeholder = '',
    error = '',
}: InputProps) => {
    return (
        <div className={`mb-4 ${className}`}>
            <label className='block text-sm font-medium text-gray-700 mb-1'>
                {label}
                {required && <span className="text-red-500">*</span>}
            </label>
            <input
                type={type}
                value={value}
                onChange={onChange}
                placeholder={placeholder}
                className={`w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    error ? 'border-red-500' : 'border-gray-300'
                }`}
            />
            {error && (
                <p className="mt-1 text-sm text-red-600">{error}</p>
            )}
        </div>
    )
}