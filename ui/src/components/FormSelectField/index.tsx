import { useEffect } from "react";
import { Control, FieldValues, Path, useFormContext } from "react-hook-form";
import { FormControl, FormField, FormItem, FormMessage } from "../ui/form";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select";


interface FormSelectProps<T extends FieldValues> {
    name: Path<T>;
    placeholder: string;
    control: Control<T>;
    options: { value: string; label: string }[];
    disabled?: boolean;
    className?: string;
    clearInvalid?: boolean;
}

const FormSelectField = <T extends FieldValues>({
    name,
    placeholder,
    control,
    options,
    disabled,
    className,
    clearInvalid
}: FormSelectProps<T>) => {
    const methods = useFormContext()

    return <FormField
        control={control}
        name={name}
        render={({ field }) => {

            useEffect(() => {
                if (!options.length) return methods.resetField<string>(name, { defaultValue: "" });
                if (!clearInvalid) return;
                const isValidOption = options.some(option => option.value === field.value);
                if (!isValidOption) methods.resetField<string>(name, { defaultValue: "" });
            }, [options, field.value, clearInvalid]);

            return (
                <FormItem>
                    <FormControl>
                        <Select name={name} onValueChange={field.onChange} defaultValue={field.value} value={field.value} disabled={disabled}>
                            <SelectTrigger className={className}>
                                <SelectValue placeholder={placeholder} />
                            </SelectTrigger>
                            <SelectContent>
                                {options.map((option) => (
                                    <SelectItem key={option.value} value={option.value}>
                                        {option.label}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </FormControl>
                    <FormMessage />
                </FormItem>
            )
        }}
    />
}

export default FormSelectField;