import { ReactNode } from "react";

type ButtonProps = {
    children: ReactNode;
    type?: 'button' | 'submit' | 'reset';
    variant?: 'primary' | 'secondary';
    disabled?: boolean;
    className?: string;
    onClick?: () => void;
};

export const Button = ({
    children,
    type = 'button',
    variant = 'primary',
    disabled = false,
    className = '',
    ...props
}: ButtonProps) => {
    const baseStyle = 'px-4 py-2 rounded font-medium';
    const variants = {
        primary: "bg-blue-600 text-white hover:bg-blue-700",
        secondary: "bg-gray-200 text-gray-800 hover:bg-gray-300",
    };

    return (
        <button
        type={type}
        className={`${baseStyle} ${variants[variant]} ${className}`}
        disabled={disabled}
        {...props}
        >
            {children}
        </button>
    )
}